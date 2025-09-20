package friend

import (
	"context"

	"github.com/wujunhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/types"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendDeleteLogic {
	return &FriendDeleteLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *FriendDeleteLogic) FriendDelete(req *types.FriendDeleteReq) (resp *types.FriendDeleteResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	_, err = l.svcCtx.Social.FriendDelete(l.ctx, &socialclient.FriendDeleteReq{UserId: uid, FriendUid: req.FriendUid})
	return &types.FriendDeleteResp{}, err
}
