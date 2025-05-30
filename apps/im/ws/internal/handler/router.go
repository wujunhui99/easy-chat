package handler

import (
	"github.com/wujunhui99/easy-chat/apps/im/ws/internal/handler/conversation"
	"github.com/wujunhui99/easy-chat/apps/im/ws/internal/handler/push"
	"github.com/wujunhui99/easy-chat/apps/im/ws/internal/handler/user"
	"github.com/wujunhui99/easy-chat/apps/im/ws/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/im/ws/websocket"
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
		{
			Method:  "conversation.markChat",
			Handler: conversation.MarkRead(svc),
		},
	})

}
