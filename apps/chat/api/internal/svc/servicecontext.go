package svc

import (
	"github.com/wujunhui99/easy-chat/apps/chat/api/internal/config"
	"github.com/wujunhui99/easy-chat/apps/chat/api/internal/middleware"
	"github.com/wujunhui99/easy-chat/apps/chat/rpc/chatclient"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	chatclient.Chat
	userclient.User
	socialclient.Social
	*redis.Redis
	JwtParse rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Chat:       chatclient.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		User:     userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social:   socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		Redis:    redis.MustNewRedis(c.JwtTable),
		JwtParse: middleware.NewJwtParseMiddleware(c).Handle,
	}
}
