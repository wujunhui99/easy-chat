package user

import (
	websocketx "github.com/wujunhui99/easy-chat/apps/msg/msggateway/websocket"
	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/internal/svc"
)

func OnLine(svc *svc.ServiceContext) websocketx.HandlerFunc {
	return func(srv *websocketx.Server, conn *websocketx.Conn, message *websocketx.Message) {
		uids := srv.GetUsers()
		u := srv.GetUsers(conn)
		err := srv.Send(websocketx.NewMessage(u[0], uids), conn)
		srv.Info("err ", err)
	}
}
