package ws

import "github.com/junhui99/easy-chat/pkg/constants"

type (
	Msg struct {
		constants.MType `mapstructure:"msgType"`
		Content         string `mapstructure:"content"`
	}
)

type (
	Chat struct {
		ConversationId string             `mapstructure:"conversationId"`
		SendId             string                    `mapstructure:"sendId"`    
		RecvId             string                    `mapstructure:"recvId"`   
		SendTime           int64                     `mapstructure:"sendTime"`  
		ChatType       constants.ChatType `mapstructure:"chatType"`
		Msg            `mapstructure:"msg"`
	}
)
