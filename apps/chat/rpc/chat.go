package main

import (
	"flag"
	"fmt"

	"github.com/wujunhui99/easy-chat/apps/chat/rpc/chat"
	"github.com/wujunhui99/easy-chat/apps/chat/rpc/internal/config"
	"github.com/wujunhui99/easy-chat/apps/chat/rpc/internal/server"
	"github.com/wujunhui99/easy-chat/apps/chat/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/pkg/interceptor/rpcserver"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// conf.MustLoad(*configFile, &c)
	conf.LoadConfig(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		chat.RegisterChatServer(grpcServer, server.NewChatServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	defer s.Stop()

	fmt.Printf("Starting im rpc server at %s...\n", c.ListenOn)
	s.Start()
}
