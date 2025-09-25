package main

import (
	"flag"
	"fmt"

	"github.com/wujunhui99/easy-chat/apps/msg/msgtransfer/internal/config"
	"github.com/wujunhui99/easy-chat/apps/msg/msgtransfer/internal/handler"
	"github.com/wujunhui99/easy-chat/apps/msg/msgtransfer/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/msgtransfer.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.LoadConfig(*configFile, &c, conf.UseEnv())

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting message transfer service...")
	serviceGroup.Start()
}
