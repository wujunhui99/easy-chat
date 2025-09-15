package user

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/types"
	userrpc "github.com/wujunhui99/easy-chat/apps/user/rpc/user"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type KickDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 踢下其它设备(仅主设备mobile)
func NewKickDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickDeviceLogic {
	return &KickDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KickDeviceLogic) KickDevice(req *types.KickDeviceReq) (resp *types.KickDeviceResp, err error) {
	// 上下文获取当前登录用户与当前设备类型
	uid := ctxdata.GetUid(l.ctx)
	callerDev := ctxdata.GetDevicetype(l.ctx)
	if uid == "" || callerDev == "" {
		return nil, errors.WithStack(xerr.New(xerr.UNAUTHORIZED_ERROR, "缺少登录上下文"))
	}
	if !constants.IsPrimaryDevice(callerDev) {
		return nil, errors.WithStack(xerr.New(xerr.NO_PERMISSION, xerr.ErrMsg(xerr.NO_PERMISSION)))
	}
	target := strings.TrimSpace(req.TargetDeviceType)
	if target == "" {
		return nil, errors.WithStack(xerr.New(xerr.REQUEST_PARAM_ERROR, "targetDeviceType 不能为空"))
	}
	target = strings.ToLower(target)
	if !constants.IsAllowedDevice(target) {
		return nil, errors.WithStack(xerr.New(xerr.INVALID_DEVICE_TYPE, xerr.ErrMsg(xerr.INVALID_DEVICE_TYPE)))
	}
	if target == callerDev {
		return nil, errors.WithStack(xerr.New(xerr.REQUEST_PARAM_ERROR, "不能踢当前设备本身"))
	}

	r, err := l.svcCtx.User.KickDevice(l.ctx, &userrpc.KickDeviceReq{
		TargetDeviceType: target,
		CallerUid:        uid,
		CallerDeviceType: callerDev,
	})
	if err != nil {
		return nil, err
	}
	return &types.KickDeviceResp{Success: r.Success}, nil
}
