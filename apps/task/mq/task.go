/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package main

import (
	"flag"
	"fmt"

	"github.com/junhui99/easy-chat/apps/task/mq/internal/config"
	"github.com/junhui99/easy-chat/apps/task/mq/internal/handler"
	"github.com/junhui99/easy-chat/apps/task/mq/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mqueue at ...")
	serviceGroup.Start()
}
