package friend

import (
	"context"
	"strings"

	"github.com/wujunhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/types"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/userclient"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
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
	currentUserId := ctxdata.GetUid(l.ctx)

	// 先获取好友申请记录，确定申请人ID
	friendReqList, err := l.svcCtx.Social.FriendPutInList(l.ctx, &socialclient.FriendPutInListReq{
		UserId:    currentUserId,
		Direction: 1, // 收到的申请
	})
	if err != nil {
		return nil, err
	}

	// 查找对应的申请记录
	var reqUid string
	for _, friendReq := range friendReqList.List {
		if friendReq.Id == req.FriendReqId {
			reqUid = friendReq.ReqUid
			break
		}
	}

	remark := strings.TrimSpace(req.Remark)
	if reqUid != "" {
		userInfo, err := l.svcCtx.User.GetUserInfo(l.ctx, &userclient.GetUserInfoReq{Id: reqUid})
		if err == nil && userInfo.User != nil {
			if remark == "" {
				remark = userInfo.User.Nickname
			}
		}
	}

	_, err = l.svcCtx.Social.FriendPutInHandle(l.ctx, &socialclient.FriendPutInHandleReq{
		FriendReqId:  req.FriendReqId,
		UserId:       currentUserId,
		HandleResult: req.HandleResult,
		Remark:       remark, // 使用申请人的nickname作为remark
	})

	return
}
