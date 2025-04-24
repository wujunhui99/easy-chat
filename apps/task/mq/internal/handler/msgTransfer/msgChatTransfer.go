package msgTransfer

import (
	"context"

	"encoding/json"
	"fmt"

	"github.com/junhui99/easy-chat/apps/im/immodels"
	"github.com/junhui99/easy-chat/apps/im/ws/ws"
	"github.com/junhui99/easy-chat/apps/task/mq/internal/svc"
	"github.com/junhui99/easy-chat/apps/task/mq/mq"
	"github.com/junhui99/easy-chat/pkg/bitmap"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

//实现kafka消费者

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key: ", key, "value: ", value)
	var (
		data  mq.MsgChatTransfer
		msgId = primitive.NewObjectID()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	//记录数据
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err
	}
	fmt.Println("msg id",msgId)
	fmt.Println("msg daata id", data.MsgId)
	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		MsgId:          msgId.Hex(),
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	//记录消息
	chatLog := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgFrom:        0,
		MsgType:        data.MType,
		ChatType:       data.ChatType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}
	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecords.Export()
	//更新会话
	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
