package conversation

import (
	"time"

	"github.com/junhui99/easy-chat/apps/im/ws/internal/svc"
	"github.com/junhui99/easy-chat/apps/im/ws/websocket"
	"github.com/junhui99/easy-chat/apps/im/ws/ws"
	"github.com/junhui99/easy-chat/apps/task/mq/mq"
	"github.com/junhui99/easy-chat/pkg/constants"
	"github.com/junhui99/easy-chat/pkg/wuid"
	"github.com/mitchellh/mapstructure"
)

func Chat(svcCtx *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, message *websocket.Message) {
		// TODO: implement
		var data ws.Chat
		if err := mapstructure.Decode(message.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.SingleChatType:
				data.ConversationId = wuid.CombineId(srv.GetUsers(conn)[0], data.RecvId)
			case constants.GroupChatType:
				data.ConversationId = data.RecvId
			}
		}

		err := svcCtx.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			MType:          data.MType,
			Content:        data.Content,
			ChatType:       data.ChatType,
			SendTime:       time.Now().UnixNano(),
		})

		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return

		}

	}
}


func MarkRead(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 已读未读处理
		var data ws.MarkRead
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		err := svc.MsgReadTransferClient.Push(&mq.MsgMarkRead{
			ChatType:       data.ChatType,
			ConversationId: data.ConversationId,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			MsgIds:         data.MsgIds,
		})
		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}
