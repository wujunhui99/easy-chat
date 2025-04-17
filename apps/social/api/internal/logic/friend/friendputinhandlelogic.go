package friend

import (
	"context"

	"github.com/junhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/junhui99/easy-chat/apps/social/api/internal/types"
	"github.com/junhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/junhui99/easy-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (resp *types.FriendPutInHandleResp, err error) {
	// todo: add your logic here and delete this line

	_, err = l.svcCtx.Social.FriendPutInHandle(l.ctx, &socialclient.FriendPutInHandleReq{
		FriendReqId:  req.FriendReqId,
		UserId:       ctxdata.GetUid(l.ctx),
		HandleResult: req.HandleResult,
	})

	return
}
