package msgTransfer

import (
	"context"
	"fmt"

	"github.com/junhui99/easy-chat/apps/task/mq/internal/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMsgChatTransfer(svcCtx *svc.ServiceContext) kq.ConsumeHandler {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}
func (m *MsgChatTransfer) Consume(key, value string) error {
	fmt.Println("key:", key, "value:", value)
	return nil
}
