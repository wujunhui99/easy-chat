package logic

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/user"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type KickDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickDeviceLogic {
	return &KickDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KickDeviceLogic) KickDevice(in *user.KickDeviceReq) (*user.KickDeviceResp, error) {
	uid := strings.TrimSpace(in.CallerUid)
	callerDev := strings.TrimSpace(in.CallerDeviceType)
	if uid == "" || callerDev == "" {
		return nil, errors.WithStack(xerr.New(xerr.UNAUTHORIZED_ERROR, "缺少登录上下文"))
	}
	if !constants.IsPrimaryDevice(callerDev) {
		return nil, errors.WithStack(xerr.New(xerr.NO_PERMISSION, xerr.ErrMsg(xerr.NO_PERMISSION)))
	}
	if in.TargetDeviceType == "" {
		return nil, errors.WithStack(xerr.New(xerr.REQUEST_PARAM_ERROR, "目标设备类型不能为空"))
	}
	target := strings.ToLower(in.TargetDeviceType)
	if !constants.IsAllowedDevice(target) {
		return nil, errors.WithStack(xerr.New(xerr.INVALID_DEVICE_TYPE, xerr.ErrMsg(xerr.INVALID_DEVICE_TYPE)))
	}
	if target == callerDev {
		return nil, errors.WithStack(xerr.New(xerr.REQUEST_PARAM_ERROR, "不能踢当前设备"))
	}
	// 幂等删除
	_, err := l.svcCtx.Redis.Del(uid + ":" + target)
	if err != nil {
		return nil, errors.WithStack(xerr.NewDBErr())
	}
	return &user.KickDeviceResp{Success: true}, nil
}
