package logic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/junhui99/easy-chat/apps/im/rpc/imclient"
	"github.com/junhui99/easy-chat/pkg/ctxdata"

	"github.com/junhui99/easy-chat/apps/im/api/internal/svc"
	"github.com/junhui99/easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationsLogic) GetConversations(req *types.GetConversationsReq) (resp *types.GetConversationsResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	data, err := l.svcCtx.GetConversations(l.ctx, &imclient.GetConversationsReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	var res types.GetConversationsResp
	copier.Copy(&res, &data)

	return &res, err
}
