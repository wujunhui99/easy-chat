package main

import (
	"flag"
	"fmt"

	"github.com/junhui99/easy-chat/apps/im/ws/internal/config"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/handler"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/svc"
	"github.com/junhui99/easy-chat/apps/im/ws/websocket"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.LoadConfig(*configFile, &c, conf.UseEnv())
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerAck(websocket.NoAck),
	// websocket.WithServerMaxConnectionIdle(10*time.Second)
	)
	defer srv.Stop()
	// ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(srv, ctx)

	fmt.Println("Starting server at", c.ListenOn)
	srv.Start()

}
