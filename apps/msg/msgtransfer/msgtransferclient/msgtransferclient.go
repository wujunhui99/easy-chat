package msgtransferclient

import (
	"context"
	"encoding/json"


	"github.com/wujunhui99/easy-chat/apps/msg/msgtransfer/msgtransfer"
	"github.com/zeromicro/go-queue/kq"
)

type MsgChatTransfer interface {
	Push(msg *MsgChatTransfer) error
}

type MsgChatTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgChatTransferClient(addrs []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return MsgChatTransferClient{
		pusher: kq.NewPusher(addrs, topic, opts...),
	}
}
func (c *MsgChatTransferClient) Push(msg *msgtransfer.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.pusher.Push(context.Background(), string(body))
}

type MsgReadTransferClient interface {
	Push(msg *msgtransfer.MsgMarkRead) error
}
type msgReadTransferClient struct {
	pusher *kq.Pusher
}

func (m *msgReadTransferClient) Push(msg *msgtransfer.MsgMarkRead) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return m.pusher.Push(context.Background(), string(body))
}

func NewMsgReadTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgReadTransferClient {
	return &msgReadTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}
