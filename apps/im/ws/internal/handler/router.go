package handler

import (
	"github.com/junhui99/easy-chat/apps/im/ws/internal/handler/conversation"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/handler/push"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/handler/user"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/svc"
	"github.com/junhui99/easy-chat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {

	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
		{
			Method:  "push",
			Handler: push.Push(svc),
		},
	})

}
