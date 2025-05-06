package main

import (
	"flag"
	"fmt"

	"github.com/wujunhui99/easy-chat/apps/task/mq/internal/config"
	"github.com/wujunhui99/easy-chat/apps/task/mq/internal/handler"
	"github.com/wujunhui99/easy-chat/apps/task/mq/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.LoadConfig(*configFile, &c, conf.UseEnv())

	// 在 main.go 中加载完配置后，添加这行代码进行调试

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	fmt.Printf("Redis Host after config loading: %s\n", c.Redisx.Host)
	// fmt.Println(c)

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting server at", c.ListenOn)
	serviceGroup.Start()
}
