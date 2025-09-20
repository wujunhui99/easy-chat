package main

import (
	"flag"
	"fmt"

	"github.com/wujunhui99/easy-chat/apps/social/api/internal/config"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/handler"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/wujunhui99/easy-chat/pkg/resultx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
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

	// 统一返回格式（与 user api 保持一致）
	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting social api server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
