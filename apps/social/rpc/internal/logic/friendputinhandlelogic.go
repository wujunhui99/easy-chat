package logic

import (
	"context"
	"database/sql"
	"strings"

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

		nickCache := make(map[string]string)
		lookupNickname := func(id string) string {
			if v, ok := nickCache[id]; ok {
				return v
			}
			info, err := l.svcCtx.UsersModel.FindOne(ctx, id)
			if err != nil {
				l.Logger.Errorf("find nickname failed id=%s err=%v", id, err)
				return ""
			}
			nickCache[id] = info.Nickname
			return info.Nickname
		}

		buildRemark := func(holderId, friendId string) sql.NullString {
			remarkText := ""
			if holderId == firendReq.UserId {
				remarkText = strings.TrimSpace(in.Remark)
			}
			if remarkText == "" {
				remarkText = lookupNickname(friendId)
			}
			if remarkText == "" {
				return sql.NullString{}
			}
			return sql.NullString{String: remarkText, Valid: true}
		}

		for _, d := range dirs {
			remarkValue := buildRemark(d.U, d.F)

			row, ferr := l.svcCtx.FriendsModel.FindOneByUserIdFriendUid(ctx, d.U, d.F)
			if ferr != nil {
				if ferr == socialmodels.ErrNotFound { // 新建方向
					toInsert = append(toInsert, &socialmodels.Friends{
						UserId:    d.U,
						FriendUid: d.F,
						Remark:    remarkValue,
						Status:    0,
					})
					continue
				}
				return errors.Wrapf(xerr.NewDBErr(), "find friend %s->%s err %v", d.U, d.F, ferr)
			}

			if row.Status != 0 {
				row.Status = 0
				if remarkValue.Valid {
					row.Remark = remarkValue
				}
				if err := l.svcCtx.FriendsModel.Update(ctx, row); err != nil {
					return errors.Wrapf(xerr.NewDBErr(), "restore friend %s->%s err %v", d.U, d.F, err)
				}
			} else if remarkValue.Valid && (!row.Remark.Valid || row.Remark.String == "") {
				row.Remark = remarkValue
				if err := l.svcCtx.FriendsModel.Update(ctx, row); err != nil {
					return errors.Wrapf(xerr.NewDBErr(), "update friend remark %s->%s err %v", d.U, d.F, err)
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
