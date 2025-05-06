package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ImRpc     zrpc.RpcClientConf
	UserRpc   zrpc.RpcClientConf
	SocialRpc zrpc.RpcClientConf
	JwtTable  redis.RedisConf
	JwtAuth   struct {
		AccessSecret string
	}
}
