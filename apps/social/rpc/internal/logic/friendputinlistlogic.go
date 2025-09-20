package logic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/social"
	"github.com/wujunhui99/easy-chat/apps/social/socialmodels"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInListLogic) FriendPutInList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	// 兼容：未传或非法 direction 默认 1 (收到的)
	direction := in.Direction
	if direction != 2 { // 只识别 2，其余归 1
		direction = 1
	}

	var (
		reqList []*socialmodels.FriendRequests
		err     error
	)

	switch direction {
	case 2: // 我发出的 (我加别人) => req_uid = me
		reqList, err = l.svcCtx.FriendRequestsModel.ListByReqUid(l.ctx, in.UserId)
	default: // 1 收到的 (别人加我) => user_id = me
		reqList, err = l.svcCtx.FriendRequestsModel.ListByUserId(l.ctx, in.UserId)
	}
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list friend requests err %v userId %s direction %d", err, in.UserId, direction)
	}

	// 拷贝到 proto
	var respList []*social.FriendRequests
	if len(reqList) > 0 {
		copier.Copy(&respList, &reqList)
	}

	return &social.FriendPutInListResp{List: respList}, nil
}
