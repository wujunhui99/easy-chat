package friend

import (
	"context"

	"github.com/wujunhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/types"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type MutualFriendCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 共同好友数量（Redis 求交集）
func NewMutualFriendCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MutualFriendCountLogic {
	return &MutualFriendCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MutualFriendCountLogic) MutualFriendCount(req *types.MutualFriendCountReq) (resp *types.MutualFriendCountResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	r, err := l.svcCtx.Social.MutualFriendCount(l.ctx, &socialclient.MutualFriendCountReq{
		UserId:  uid,
		OtherId: req.OtherId,
	})
	if err != nil {
		return nil, err
	}
	return &types.MutualFriendCountResp{Count: r.Count}, nil
}
