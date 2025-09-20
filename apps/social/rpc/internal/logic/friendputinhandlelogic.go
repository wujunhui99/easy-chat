package logic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/social"
	"github.com/wujunhui99/easy-chat/apps/social/socialmodels"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrFriendReqBeforePass   = xerr.NewMsg("好友申请并已经通过")
	ErrFriendReqBeforeRefuse = xerr.NewMsg("好友申请已经被拒绝")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// 获取好友申请记录
	firendReq, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, uint64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by friendReqid err %v req %v ", err, in.FriendReqId)
	}

	// 已处理校验
	switch constants.HandlerResult(firendReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforeRefuse)
	}

	firendReq.HandleResult.Int64 = int64(in.HandleResult)

	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新申请
		if err := l.svcCtx.FriendRequestsModel.Update(ctx, firendReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request err %v, req %v", err, firendReq)
		}

		// 非通过无需建立关系
		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		type pair struct{ U, F string }
		dirs := []pair{{firendReq.UserId, firendReq.ReqUid}, {firendReq.ReqUid, firendReq.UserId}}
		var toInsert []*socialmodels.Friends

		for _, d := range dirs {
			// 查是否已有方向
			row, ferr := l.svcCtx.FriendsModel.FindOneByUserIdFriendUid(ctx, d.U, d.F)
			if ferr != nil {
				if ferr == socialmodels.ErrNotFound { // 插入
					toInsert = append(toInsert, &socialmodels.Friends{UserId: d.U, FriendUid: d.F, Status: 0})
					continue
				}
				return errors.Wrapf(xerr.NewDBErr(), "find friend %s->%s err %v", d.U, d.F, ferr)
			}
			// 已存在但 status!=0 (删除/拉黑/免打扰) -> 恢复为正常(是否覆盖拉黑按业务策略，这里统一恢复除非想保留拉黑需判断)
			if row.Status != 0 {
				row.Status = 0
				if err := l.svcCtx.FriendsModel.Update(ctx, row); err != nil {
					return errors.Wrapf(xerr.NewDBErr(), "restore friend %s->%s err %v", d.U, d.F, err)
				}
			}
		}

		if len(toInsert) > 0 {
			for _, f := range toInsert { // goctl 生成的 Insert 是单条，循环调用
				if _, ierr := l.svcCtx.FriendsModel.Insert(ctx, f); ierr != nil {
					return errors.Wrapf(xerr.NewDBErr(), "insert friend %v err %v", f, ierr)
				}
			}
		}
		return nil
	})

	return &social.FriendPutInHandleResp{}, err
}
