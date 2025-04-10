package conversation

import (
	"context"

	"github.com/junhui99/easy-chat/apps/im/ws/internal/logic"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/svc"
	"github.com/junhui99/easy-chat/apps/im/ws/websocket"
	"github.com/junhui99/easy-chat/apps/im/ws/ws"
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
		l := logic.NewConversation(context.Background(), srv, svcCtx)
		if err := l.Chat(&data, srv.GetUsers(conn)[0]); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		uid := srv.GetUsers(conn)[0]
		err := srv.SendByUserId(websocket.NewMessage(uid, message.Data), data.RecvId)
		srv.Info(err)

	}
}
