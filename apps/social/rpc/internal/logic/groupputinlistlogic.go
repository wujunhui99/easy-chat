package logic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/social"
	"github.com/wujunhui99/easy-chat/apps/social/socialmodels"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutinListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinListLogic {
	return &GroupPutinListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutinListLogic) GroupPutinList(in *social.GroupPutinListReq) (*social.GroupPutinListResp, error) {
	// todo: add your logic here and delete this line
	// 权限：仅群主或管理员可查看
	mem, err := l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.OperatorUid, in.GroupId)
	if err == socialmodels.ErrNotFound {
		// 非群成员归类为“无权限”
		return nil, errors.WithStack(xerr.New(xerr.NO_PERMISSION, "仅群主或管理员可查看入群申请列表"))
	}
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member gid=%s uid=%s err: %v", in.GroupId, in.OperatorUid, err)
	}
	role := constants.GroupRoleLevel(mem.RoleLevel)
	if !(role == constants.CreatorGroupRoleLevel || role == constants.ManagerGroupRoleLevel) {
		return nil, errors.WithStack(xerr.New(xerr.NO_PERMISSION, "仅群主或管理员可查看入群申请列表"))
	}

	groupReqs, err := l.svcCtx.GroupRequestsModel.ListNoHandler(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group req err %v req %v", err, in.GroupId)
	}

	var respList []*social.GroupRequests
	copier.Copy(&respList, groupReqs)

	return &social.GroupPutinListResp{
		List: respList,
	}, nil
}
