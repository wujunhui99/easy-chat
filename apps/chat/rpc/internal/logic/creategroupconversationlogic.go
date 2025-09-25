package logic

import (
	"context"

	"github.com/pkg/errors"
	chatmodels "github.com/wujunhui99/easy-chat/apps/chat/chatmodels"
	"github.com/wujunhui99/easy-chat/apps/chat/rpc/chat"
	"github.com/wujunhui99/easy-chat/apps/chat/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupConversationLogic) CreateGroupConversation(in *chat.CreateGroupConversationReq) (*chat.CreateGroupConversationResp, error) {
	// todo: add your logic here and delete this line
	res := &chat.CreateGroupConversationResp{}
	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	if err == nil {
		return res, nil
	}
	if err != chatmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.FindOne err %v, req %v", err, in.GroupId)
	}
	err = l.svcCtx.ConversationModel.Insert(l.ctx, &chatmodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.Insert err %v, req %v", err, in.GroupId)
	}
	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&chat.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})
	return res, err

}
