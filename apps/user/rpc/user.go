package main

import (
	"flag"
	"fmt"

	"github.com/wujunhui99/easy-chat/apps/user/rpc/internal/config"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/internal/server"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/user"
	"github.com/wujunhui99/easy-chat/pkg/interceptor/rpcserver"
	"github.com/wujunhui99/easy-chat/pkg/wuid"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// conf.MustLoad(*configFile, &c,conf.UseEnv())
	conf.LoadConfig(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	if err := ctx.SetRootToken(); err != nil {
		panic(err)
	}
	wuid.Init(c.Mysql.DataSource)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	defer s.Stop()
	fmt.Printf("Starting user rpc server at %s...\n", c.ListenOn)
	s.Start()
}
