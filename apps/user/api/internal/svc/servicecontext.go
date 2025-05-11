package svc

import (
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/config"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/middleware"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	userclient.User
	*redis.Redis
	JwtParseMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		User:               userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Redis:              redis.MustNewRedis(c.JwtTable),
		JwtParseMiddleware: middleware.NewJwtParseMiddleware(c).Handle,
	}
}
