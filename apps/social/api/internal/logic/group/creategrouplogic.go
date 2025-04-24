package group

import (
	"context"

	"github.com/junhui99/easy-chat/apps/im/rpc/imclient"
	"github.com/junhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/junhui99/easy-chat/apps/social/api/internal/types"
	"github.com/junhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/junhui99/easy-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUid(l.ctx)

	// 创建群
	res, err := l.svcCtx.Social.GroupCreate(l.ctx, &socialclient.GroupCreateReq{
		Name:       req.Name,
		Icon:       req.Icon,
		CreatorUid: uid,
	})
	if err != nil {
		return nil, err
	}

	if res.Id == "" {
		return nil, err
	}

	_, err = l.svcCtx.Im.CreateGroupConversation(l.ctx, &imclient.CreateGroupConversationReq{
		GroupId:  res.Id,
		CreateId: uid,
	})

	return
}
