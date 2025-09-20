package group

import (
	"context"

	"github.com/wujunhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/types"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type MyCreatedGroupsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 我创建的群列表
func NewMyCreatedGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyCreatedGroupsLogic {
	return &MyCreatedGroupsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MyCreatedGroupsLogic) MyCreatedGroups(req *types.MyCreatedGroupsReq) (resp *types.MyCreatedGroupsResp, err error) {
	uid := ctxdata.GetUid(l.ctx)

	r, err := l.svcCtx.Social.MyCreatedGroups(l.ctx, &socialclient.MyCreatedGroupsReq{UserId: uid})
	if err != nil {
		return nil, err
	}

	// 映射到 API types
	out := &types.MyCreatedGroupsResp{List: make([]*types.Groups, 0, len(r.List))}
	for _, g := range r.List {
		out.List = append(out.List, &types.Groups{
			Id:              g.Id,
			Name:            g.Name,
			Icon:            g.Icon,
			Status:          int64(g.Status),
			GroupType:       int64(g.GroupType),
			IsVerify:        g.IsVerify,
			Notification:    g.Notification,
			NotificationUid: g.NotificationUid,
			CreatedAt:       g.CreatedAt,
		})
	}
	return out, nil
}
