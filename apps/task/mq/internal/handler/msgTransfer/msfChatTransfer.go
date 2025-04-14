package msgTransfer

import (
	"context"
	"encoding/json"

	"github.com/junhui99/easy-chat/apps/im/immodels"
	"github.com/junhui99/easy-chat/apps/im/ws/websocket"
	"github.com/junhui99/easy-chat/apps/task/mq/internal/svc"
	"github.com/junhui99/easy-chat/apps/task/mq/mq"
	"github.com/junhui99/easy-chat/pkg/constants"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMsgChatTransfer(svcCtx *svc.ServiceContext) kq.ConsumeHandler {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (m *MsgChatTransfer) Consume(key, value string) error {
	var (
		data mq.MsgChatTransfer
		ctx  = context.Background()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}
	if err := m.addChatLog(ctx, data); err != nil {
		return err
	}
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameNoAck,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})


	return nil
}
func (m *MsgChatTransfer) addChatLog(ctx context.Context, data mq.MsgChatTransfer) error {
	// 1.添加聊天记录
	chatLog := &immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgType :          data.MType,
		MsgContent: data.Content,
		ChatType: 	 data.ChatType,
	}
	if err := m.svcCtx.ChatLogModel.Insert(ctx, chatLog); err != nil {
		return err
	}

	return nil
}
