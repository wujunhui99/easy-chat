package logic

import (
	"context"

	"github.com/wujunhui99/easy-chat/apps/social/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/social"
	"github.com/wujunhui99/easy-chat/apps/social/socialmodels"

	"github.com/zeromicro/go-zero/core/logx"
)

type MyCreatedGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMyCreatedGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyCreatedGroupsLogic {
	return &MyCreatedGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MyCreatedGroupsLogic) MyCreatedGroups(in *social.MyCreatedGroupsReq) (*social.MyCreatedGroupsResp, error) {
	// 查询我创建的群
	list, err := l.svcCtx.GroupsModel.ListByCreatorUid(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	// 映射到 proto 类型
	var ret []*social.Groups
	ret = make([]*social.Groups, 0, len(list))
	for _, g := range list {
		ret = append(ret, modelGroupToProto(g))
	}

	return &social.MyCreatedGroupsResp{List: ret}, nil
}

// modelGroupToProto 将 socialmodels.Groups 转为 rpc social.Groups
func modelGroupToProto(g *socialmodels.Groups) *social.Groups {
	if g == nil {
		return nil
	}
	return &social.Groups{
		Id:              g.Id,
		Name:            g.Name,
		Icon:            g.Icon,
		Status:          int32(g.Status.Int64),
		CreatorUid:      g.CreatorUid,
		GroupType:       int32(g.GroupType),
		IsVerify:        g.IsVerify,
		Notification:    g.Notification.String,
		NotificationUid: g.NotificationUid.String,
		CreatedAt:       g.CreatedAt.Unix(),
	}
}
