package main

import (
	"flag"
	"fmt"

	"github.com/junhui99/easy-chat/apps/social/api/internal/config"
	"github.com/junhui99/easy-chat/apps/social/api/internal/handler"
	"github.com/junhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/junhui99/easy-chat/middleware"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// conf.MustLoad(*configFile, &c)
	conf.LoadConfig(*configFile, &c, conf.UseEnv())
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	tokenRedisCheck := middleware.NewRedisTokenCheckMiddleware(ctx.Redis)
	server.Use(tokenRedisCheck.Handle)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
