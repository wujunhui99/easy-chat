package mqclient

import (
	"encoding/json"

	"github.com/junhui99/easy-chat/apps/task/mq/mq"
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
func (c *MsgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.pusher.Push(string(body))
}
