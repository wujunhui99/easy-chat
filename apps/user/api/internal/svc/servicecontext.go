package svc

import (
	"github.com/junhui99/easy-chat/apps/user/api/internal/config"
	"github.com/junhui99/easy-chat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
