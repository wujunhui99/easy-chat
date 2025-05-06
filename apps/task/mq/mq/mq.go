package mq

import "github.com/wujunhui99/easy-chat/pkg/constants"

type MsgChatTransfer struct {
	MsgId              string `mapstructure:"msgId"`
	ConversationId     string `json:"conversationId"` // 聊天会话的唯一标识符
	constants.ChatType `json:"chatType"`
	SendId             string   `json:"sendId"` // 发送者的唯一标识符
	RecvId             string   `json:"recvId"` // 接收者的唯一标识符
	RecvIds            []string `json:"recvIds"`
	SendTime           int64    `json:"sendTime"` // 消息发送的时间戳
	constants.MType    `json:"mType"`
	Content            string `json:"content"` // 消息的实际内容
}

type MsgMarkRead struct {
	constants.ChatType `json:"chatType"`
	ConversationId     string   `json:"conversationId"`
	SendId             string   `json:"sendId"`
	RecvId             string   `json:"recvId"`
	MsgIds             []string `json:"msgIds"`
}
