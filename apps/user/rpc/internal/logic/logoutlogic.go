package logic

import (
	"context"
	"errors"

	"github.com/wujunhui99/easy-chat/apps/user/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *user.LogoutReq) (*user.LogoutResp, error) {
	if in.Id == "" || in.DeviceType == "" {
		return nil, errors.New("缺少必需字段")
	}
	cnt, err := l.svcCtx.Redis.Del(in.Id + ":" + in.DeviceType)
	if err != nil {
		return nil, err
	}
	if cnt == 0 {
		return nil, errors.New("未找到登录状态")
	}
	return &user.LogoutResp{Success: 1}, nil
}
