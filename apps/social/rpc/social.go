package main

import (
	"flag"
	"fmt"

	"im-zero/east-chat/apps/social/rpc/internal/config"
	"im-zero/east-chat/apps/social/rpc/internal/server"
	"im-zero/east-chat/apps/social/rpc/internal/svc"
	"im-zero/east-chat/apps/social/rpc/social"
	"im-zero/east-chat/pkg/interceptor/rpcserver"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		social.RegisterSocialServer(grpcServer, server.NewSocialServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}