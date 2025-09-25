package svc

import (
	"github.com/wujunhui99/easy-chat/apps/chat/rpc/chatclient"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/config"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/middleware"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	User   userclient.User
	Social socialclient.Social
	Chat   chatclient.Chat
	*redis.Redis
	JwtParse rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		User:     userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social:   socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		Chat:     chatclient.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		Redis:    redis.MustNewRedis(c.JwtTable),
		JwtParse: middleware.NewJwtParseMiddleware(c).Handle,
	}
}
