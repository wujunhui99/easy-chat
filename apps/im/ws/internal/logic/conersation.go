package logic

import (
	"context"
	"time"

	"github.com/wujunhui99/easy-chat/apps/im/immodels"
	"github.com/wujunhui99/easy-chat/apps/im/ws/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/im/ws/websocket"
	"github.com/wujunhui99/easy-chat/apps/im/ws/ws"
	"github.com/wujunhui99/easy-chat/pkg/wuid"
)

type Conversation struct {
	ctx context.Context
	srv *websocket.Server
	svc *svc.ServiceContext
}

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

func (l *Conversation) Chat(data *ws.Chat, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,

		SendId:     userId,
		RecvId:     data.RecvId,
		SendTime:   time.Now().UnixNano(),
		MsgType:    data.MType,
		MsgContent: data.Content,
		ChatType:   data.ChatType,
	}
	err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)
	return err
}
