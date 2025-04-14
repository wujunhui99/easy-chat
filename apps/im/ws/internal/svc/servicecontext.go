package svc

import (
	"github.com/junhui99/easy-chat/apps/im/immodels"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/config"
	"github.com/junhui99/easy-chat/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config
	immodels.ChatLogModel
	MsgChatTransferClient mqclient.MsgChatTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
	}
}
