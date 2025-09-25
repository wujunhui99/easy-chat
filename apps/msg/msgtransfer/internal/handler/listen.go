package handler

import (
	"github.com/wujunhui99/easy-chat/apps/msg/msgtransfer/internal/handler/msgTransfer"
	"github.com/wujunhui99/easy-chat/apps/msg/msgtransfer/internal/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Listen struct {
	svc *svc.ServiceContext
}

func NewListen(svc *svc.ServiceContext) *Listen {
	return &Listen{svc: svc}
}

//返回多个消费者对象

func (l *Listen) Services() []service.Service {

	return []service.Service{
		kq.MustNewQueue(l.svc.Config.MsgReadTransfer, msgTransfer.NewMsgReadTransfer(l.svc)),
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer, msgTransfer.NewMsgChatTransfer(l.svc)),
	}
}
