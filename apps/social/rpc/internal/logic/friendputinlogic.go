package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/social"
	"github.com/wujunhui99/easy-chat/apps/social/socialmodels"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// 语义明确：in.UserId = target(被申请人) , in.ReqUid = requester(申请人)
	target := in.UserId
	requester := in.ReqUid

	// 不能加自己
	if target == requester {
		return &social.FriendPutInResp{}, nil
	}

	// 是否已经是好友（requester -> target）
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, requester, target)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends requester=%s target=%s err %v", requester, target, err)
	}
	if friends != nil { // 已是好友直接返回幂等
		return &social.FriendPutInResp{}, nil
	}

	// 申请记录是否已存在（同一对）
	friendReqs, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, requester, target)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendRequest requester=%s target=%s err %v", requester, target, err)
	}
	if friendReqs != nil { // 已存在未处理或正在处理中
		return &social.FriendPutInResp{}, nil
	}

	// 插入申请
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId:       target, // 被申请人
		ReqUid:       requester,
		ReqMsg:       sql.NullString{Valid: len(in.ReqMsg) > 0, String: in.ReqMsg},
		ReqTime:      time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{Int64: int64(constants.NoHandlerResult), Valid: true},
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest requester=%s target=%s err %v", requester, target, err)
	}
	return &social.FriendPutInResp{}, nil
}
