/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package handler

import (
	"github.com/junhui99/easy-chat/apps/im/ws/internal/handle/conversation"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/handle/push"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/handle/user"
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
