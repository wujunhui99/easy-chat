package group

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/junhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/junhui99/easy-chat/apps/social/api/internal/types"
	"github.com/junhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/junhui99/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupListLogic) GroupList(req *types.GroupListRep) (resp *types.GroupListResp, err error) {
	// todo: add your logic here and delete this line

	uid := ctxdata.GetUid(l.ctx)
	list, err := l.svcCtx.Social.GroupList(l.ctx, &socialclient.GroupListReq{
		UserId: uid,
	})

	if err != nil {
		return nil, err
	}

	var respList []*types.Groups
	copier.Copy(&respList, list.List)

	return &types.GroupListResp{List: respList}, nil
}
