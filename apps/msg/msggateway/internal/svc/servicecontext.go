package svc

import (
	"github.com/wujunhui99/easy-chat/apps/chat/chatmodels"
	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/internal/config"
	"github.com/wujunhui99/easy-chat/apps/msg/msgtransfer/msgtransferclient"
	"github.com/wujunhui99/easy-chat/pkg/middleware/tokenmatch"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	chatmodels.ChatLogModel
	msgtransferclient.MsgChatTransferClient
	msgtransferclient.MsgReadTransferClient
	TokenMatch *tokenmatch.TokenMatch
}

func NewServiceContext(c config.Config) *ServiceContext {
	jwtRedis := redis.MustNewRedis(c.JwtTable)
	return &ServiceContext{
		Config:                c,
		ChatLogModel:          chatmodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		MsgChatTransferClient: msgtransferclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		MsgReadTransferClient: msgtransferclient.NewMsgReadTransferClient(c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic),
		TokenMatch:            tokenmatch.New(jwtRedis, tokenmatch.Config{AccessSecret: c.JwtAuth.AccessSecret}),
	}
}
