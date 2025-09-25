package msggateway

import "github.com/wujunhui99/easy-chat/pkg/constants"

type (
	Msg struct {
		MsgId           string `mapstructure:"msgId"`
		constants.MType `mapstructure:"msgType"`
		ReadRecords     map[string]string `mapstructure:"readRecords"`
		Content         string            `mapstructure:"content"`
	}
)

type (
	Chat struct {
		ConversationId string             `mapstructure:"conversationId"`
		SendId         string             `mapstructure:"sendId"`
		RecvId         string             `mapstructure:"recvId"`
		SendTime       int64              `mapstructure:"sendTime"`
		ChatType       constants.ChatType `mapstructure:"chatType"`
		Msg            `mapstructure:"msg"`
	}
)

type Push struct {
	ChatType        constants.ChatType `mapstructure:"chatType"`
	ConversationId  string             `mapstructure:"conversationId"`
	RecvId          string             `mapstructure:"recvId"`
	RecvIds         []string           `mapstructure:"recvIds"`
	SendId          string             `mapstructure:"sendId"`
	constants.MType `mapstructure:"mType"`
	Content         string                `mapstructure:"content"`
	ReadRecords     map[string]string     `mapstructure:"readRecords"`
	MsgId           string                `mapstructure:"msgId"`
	ContentType     constants.ContentType `mapstructure:"contentType"`

	SendTime int64 `mapstructure:"sendTime"`
}

type MarkRead struct {
	constants.ChatType `mapstructure:"chatType"`
	RecvId             string   `mapstructure:"recvId"`
	ConversationId     string   `mapstructure:"conversationId"`
	MsgIds             []string `mapstructure:"msgIds"`
}
