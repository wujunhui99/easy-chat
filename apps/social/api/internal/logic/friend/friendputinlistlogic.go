package friend

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/types"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	uid := ctxdata.GetUid(l.ctx)

	direction := req.Direction
	if direction != 2 { // 只认 2，其余归1
		direction = 1
	}

	list, err := l.svcCtx.Social.FriendPutInList(l.ctx, &socialclient.FriendPutInListReq{
		UserId:    uid,
		Direction: direction,
	})
	if err != nil {
		return nil, err
	}

	var respList []*types.FriendRequests
	if len(list.List) > 0 {
		copier.Copy(&respList, list.List)
	}

	return &types.FriendPutInListResp{List: respList}, nil
}
