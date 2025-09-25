package push

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/msggateway"
	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/websocket"
	"github.com/wujunhui99/easy-chat/pkg/constants"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data msggateway.Push

		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err))
			return
		}
		fmt.Println("exe push")
		fmt.Println("recv id", data.RecvId)
		fmt.Println("msg id", data.MsgId)
		// 发送的目标
		switch data.ChatType {
		case constants.SingleChatType:
			single(srv, &data, data.RecvId)
		case constants.GroupChatType:
			// fmt.Println("group push x", data.MsgId)
			group(srv, &data)
		}
	}
}

func single(srv *websocket.Server, data *msggateway.Push, recvId string) error {
	conns := srv.GetUserConns(recvId)
	if len(conns) == 0 {
		// todo: 目标离线
		return nil
	}

	message := websocket.NewMessage(data.SendId, &msggateway.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: msggateway.Msg{
			ReadRecords: data.ReadRecords,
			MsgId:       data.MsgId,
			MType:       data.MType,
			Content:     data.Content,
		},
	})
	for _, rconn := range conns {
		if err := srv.Send(message, rconn); err != nil {
			srv.Errorf("push single send err %v uid %s device %s", err, recvId, rconn.DeviceId)
		}
	}
	return nil

}

func group(srv *websocket.Server, data *msggateway.Push) error {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(id)
	}
	return nil
}
