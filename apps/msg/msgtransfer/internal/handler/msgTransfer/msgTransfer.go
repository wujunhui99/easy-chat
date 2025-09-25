package msgTransfer

import (
	"context"

	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/msggateway"
	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/websocket"
	"github.com/wujunhui99/easy-chat/apps/msg/msgtransfer/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/zeromicro/go-zero/core/logx"
)

type baseMsgTransfer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBaseMsgTransfer(svc *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		svcCtx: svc,
		Logger: logx.WithContext(context.Background()),
	}
}

// 转发消息
func (m *baseMsgTransfer) Transfer(ctx context.Context, data *msggateway.Push) error {
	var err error
	switch data.ChatType {
	case constants.GroupChatType:
		err = m.group(ctx, data)
	case constants.SingleChatType:
		err = m.single(ctx, data)
	}
	return err
}

// 私聊转发
func (m *baseMsgTransfer) single(ctx context.Context, data *msggateway.Push) error {
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

// 群聊转发
func (m *baseMsgTransfer) group(ctx context.Context, data *msggateway.Push) error {
	// 就要查询，群的用户
	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}
	data.RecvIds = make([]string, 0, len(users.List))

	for _, members := range users.List {
		if members.UserId == data.SendId {
			continue
		}

		data.RecvIds = append(data.RecvIds, members.UserId)
	}

	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}
