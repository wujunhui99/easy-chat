package friend

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/types"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/userclient"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	// 当前登录用户 = 申请人(requester)
	requester := ctxdata.GetUid(l.ctx)
	target := req.UserId // body 中 user_uid 表示要加的那个用户

	// 自己加自己直接忽略（也可返回错误）
	if target == requester {
		return &types.FriendPutInResp{}, nil
	}

	// 校验目标用户存在（风格对齐 login：使用通用 code 携带业务文案）
	userResp, uerr := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{Ids: []string{target}})
	if uerr != nil {
		return nil, errors.WithStack(xerr.New(xerr.DB_ERROR, "查询用户失败"))
	}
	if userResp == nil || len(userResp.User) == 0 {
		return nil, errors.WithStack(xerr.New(xerr.SERVER_COMMON_ERROR, "用户不存在"))
	}

	_, err = l.svcCtx.Social.FriendPutIn(l.ctx, &socialclient.FriendPutInReq{
		UserId:  target,    // proto 中 userId 作为被申请人(目标)
		ReqUid:  requester, // proto 中 reqUid 作为申请人
		ReqMsg:  req.ReqMsg,
		ReqTime: req.ReqTime,
	})
	return
}
