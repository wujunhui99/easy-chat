package logic

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/user"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginDevicesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLoginDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginDevicesLogic {
	return &GetLoginDevicesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLoginDevicesLogic) GetLoginDevices(in *user.GetLoginDevicesReq) (*user.GetLoginDevicesResp, error) {
	uid := in.CallerUid
	if uid == "" {
		return nil, errors.WithStack(xerr.New(xerr.UNAUTHORIZED_ERROR, "缺少登录上下文"))
	}
	devices := make([]string, 0, len(constants.AllowedDevices))
	for d := range constants.AllowedDevices {
		key := uid + ":" + d
		val, err := l.svcCtx.Redis.Get(key)
		if err != nil {
			return nil, errors.WithStack(xerr.NewDBErr())
		}
		if val != "" { // 存在 token 视为该设备在线
			devices = append(devices, d)
		}
	}
	if len(devices) > 1 {
		sort.Strings(devices)
	}
	return &user.GetLoginDevicesResp{Devices: devices}, nil
}
