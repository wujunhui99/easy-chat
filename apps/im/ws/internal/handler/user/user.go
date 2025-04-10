package user

import (
	"github.com/junhui99/easy-chat/apps/im/ws/internal/svc"
	websocketx "github.com/junhui99/easy-chat/apps/im/ws/websocket"
)

func OnLine(svc *svc.ServiceContext) websocketx.HandlerFunc {
	return func(srv *websocketx.Server, conn *websocketx.Conn, message *websocketx.Message) {
		uids := srv.GetUsers()
		u := srv.GetUsers(conn)
		err := srv.Send(websocketx.NewMessage(u[0], uids), conn)
		srv.Info("err ", err)
	}
}
