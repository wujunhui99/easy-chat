package logic

import (
	"context"
	"database/sql"

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
	ErrGroupReqBeforePass   = xerr.NewMsg("请求已通过")
	ErrGroupReqBeforeRefuse = xerr.NewMsg("请求已拒绝")
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	// todo: add your logic here and delete this line

	groupReq, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, int64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend req err %v req %v", err, in.GroupReqId)
	}

	switch constants.HandlerResult(groupReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforeRefuse)
	}

	groupReq.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}

	err = l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.GroupRequestsModel.Update(l.ctx, session, groupReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend req err %v req %v", err, groupReq)
		}

		if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
			return nil
		}

		groupMember := &socialmodels.GroupMembers{
			GroupId:     groupReq.GroupId,
			UserId:      groupReq.ReqId,
			RoleLevel:   int(constants.AtLargeGroupRoleLevel),
			OperatorUid: in.HandleUid,
		}
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, groupMember)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
		}

		return nil
	})

	if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
		return &social.GroupPutInHandleResp{}, err
	}

	return &social.GroupPutInHandleResp{
		GroupId: groupReq.GroupId,
	}, err
}
