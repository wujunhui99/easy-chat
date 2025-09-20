package logic

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/social"
	"github.com/wujunhui99/easy-chat/apps/social/socialmodels"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

// 约定：删除好友为单向软删：仅将 userId->friendUid 方向 status 置为 1（删除）
// 若记录不存在则幂等返回成功；不修改对方方向
// 后续重新加好友时再恢复/创建该方向

type FriendDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendDeleteLogic {
	return &FriendDeleteLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FriendDeleteLogic) FriendDelete(in *social.FriendDeleteReq) (*social.FriendDeleteResp, error) {
	// 查找该方向
	row, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.FriendUid)
	if err != nil {
		if errors.Is(err, socialmodels.ErrNotFound) { // 幂等
			return &social.FriendDeleteResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend %s->%s err %v", in.UserId, in.FriendUid, err)
	}
	// 已是删除/拉黑等，如果 status==1 直接返回；如果是拉黑(2) 或免打扰(3) 且希望删除，可改为1
	if row.Status == 1 { // 已删除
		return &social.FriendDeleteResp{}, nil
	}
	// 仅更新为删除；保留 remark/add_source
	row.Status = 1
	if err := l.svcCtx.FriendsModel.Update(l.ctx, row); err != nil {
		// 处理可能的并发：记录被别的事务删除 -> 视为成功
		if errors.Is(err, sql.ErrNoRows) {
			return &social.FriendDeleteResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "update friend delete %s->%s err %v", in.UserId, in.FriendUid, err)
	}
	return &social.FriendDeleteResp{}, nil
}
