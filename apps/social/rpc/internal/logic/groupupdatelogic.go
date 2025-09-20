package logic

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/social"
	"github.com/wujunhui99/easy-chat/apps/social/socialmodels"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUpdateLogic {
	return &GroupUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新群信息（带字段掩码）
func (l *GroupUpdateLogic) GroupUpdate(in *social.GroupUpdateReq) (*social.GroupUpdateResp, error) {
	// 校验 updateMask
	if len(in.UpdateMask) == 0 {
		return nil, errors.WithStack(xerr.NewMsg("update_mask 不能为空"))
	}
	// 允许更新的字段集合
	allowed := map[string]struct{}{
		"name":      {},
		"icon":      {},
		"is_verify": {},
	}
	// 将 mask 统一成小写
	var fields []string
	for _, f := range in.UpdateMask {
		key := strings.ToLower(strings.TrimSpace(f))
		if _, ok := allowed[key]; !ok {
			return nil, errors.WithStack(xerr.NewMsg("update_mask 包含不支持的字段: " + f))
		}
		fields = append(fields, key)
	}

	// 读群信息
	g, err := l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group %s err: %v", in.GroupId, err)
	}

	// 权限校验：必须是创建者或管理员
	mem, err := l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.OperatorUid, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member gid=%s uid=%s err: %v", in.GroupId, in.OperatorUid, err)
	}
	role := constants.GroupRoleLevel(mem.RoleLevel)
	if !(role == constants.CreatorGroupRoleLevel || role == constants.ManagerGroupRoleLevel) {
		return nil, errors.WithStack(xerr.New(xerr.NO_PERMISSION, "无权限操作该群"))
	}

	// 应用字段
	newData := &socialmodels.Groups{
		Id:              g.Id,
		Name:            g.Name,
		Icon:            g.Icon,
		Status:          g.Status,
		CreatorUid:      g.CreatorUid,
		GroupType:       g.GroupType,
		IsVerify:        g.IsVerify,
		Notification:    g.Notification,
		NotificationUid: g.NotificationUid,
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
	}
	for _, f := range fields {
		switch f {
		case "name":
			if strings.TrimSpace(in.Name) == "" {
				return nil, errors.WithStack(xerr.NewMsg("name 不能为空"))
			}
			newData.Name = in.Name
		case "icon":
			newData.Icon = in.Icon
		case "is_verify":
			newData.IsVerify = in.IsVerify
		}
	}

	// 事务更新（ExecCtx 会处理缓存）
	// 这里不需要开启外部事务，直接 Update 即可
	if err := l.svcCtx.GroupsModel.Update(l.ctx, newData); err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "update group id=%s err: %v", in.GroupId, err)
	}

	_ = sql.ErrNoRows // silence potential unused import in some build contexts
	return &social.GroupUpdateResp{}, nil
}
