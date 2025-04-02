/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package logic

import (
	"path/filepath"

	"github.com/junhui99/easy-chat/apps/user/rpc/internal/config"
	"github.com/junhui99/easy-chat/apps/user/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
)

var svcCtx *svc.ServiceContext

func init() {
	var c config.Config
	conf.LoadConfig(filepath.Join("../../etc/dev/user.yaml"),&c,conf.UseEnv())
	svcCtx = svc.NewServiceContext(c)
}
