package group

import (
	"context"

	"github.com/wujunhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/types"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新群信息（带 update_mask）
func NewGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUpdateLogic {
	return &GroupUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUpdateLogic) GroupUpdate(req *types.GroupUpdateReq) (resp *types.GroupUpdateResp, err error) {
	operator := ctxdata.GetUid(l.ctx)
	// 直接将 API 的字段和 update_mask 透传给 RPC
	_, err = l.svcCtx.Social.GroupUpdate(l.ctx, &socialclient.GroupUpdateReq{
		GroupId:     req.GroupId,
		Name:        req.Name,
		Icon:        req.Icon,
		IsVerify:    req.IsVerify,
		UpdateMask:  req.UpdateMask,
		OperatorUid: operator,
	})
	if err != nil {
		return nil, err
	}
	return &types.GroupUpdateResp{}, nil
}
