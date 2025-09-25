package svc

import (
	"github.com/wujunhui99/easy-chat/apps/chat/chatmodels"
	"github.com/wujunhui99/easy-chat/apps/chat/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	chatmodels.ChatLogModel
	chatmodels.ConversationsModel
	chatmodels.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		ChatLogModel:       chatmodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationsModel: chatmodels.MustConversationsModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel:  chatmodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
	}
}
