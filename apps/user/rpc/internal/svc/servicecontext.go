package svc

import (
	"github.com/junhui99/easy-chat/apps/user/models"
	"github.com/junhui99/easy-chat/apps/user/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	UsersModel models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}
