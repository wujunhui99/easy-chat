package conversation

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/msggateway"
	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/websocket"
	"github.com/wujunhui99/easy-chat/apps/msg/msgtransfer/msgtransfer"

	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/wuid"
)

func Chat(svcCtx *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, message *websocket.Message) {
		// TODO: implement
		var data msggateway.Chat
		if err := mapstructure.Decode(message.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.SingleChatType:
				data.ConversationId = wuid.CombineId(conn.Uid, data.RecvId)
			case constants.GroupChatType:
				data.ConversationId = data.RecvId
			}
		}

		err := svcCtx.MsgChatTransferClient.Push(&msgtransfer.MsgChatTransfer{
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
		var data msggateway.MarkRead
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		err := svc.MsgReadTransferClient.Push(&msgtransfer.MsgMarkRead{
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
