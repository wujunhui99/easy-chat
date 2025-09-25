package main

import (
	"flag"
	"fmt"

	"github.com/wujunhui99/easy-chat/apps/im/api/internal/config"
	"github.com/wujunhui99/easy-chat/apps/im/api/internal/handler"
	"github.com/wujunhui99/easy-chat/apps/im/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// conf.MustLoad(*configFile, &c)
	conf.LoadConfig(*configFile, &c, conf.UseEnv())
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting im api server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
