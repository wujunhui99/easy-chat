package logic

import (
	"context"
	"strings"
	"time"

	"github.com/wujunhui99/easy-chat/pkg/constants"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/user/models"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/user"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/wujunhui99/easy-chat/pkg/encrypt"
	"github.com/wujunhui99/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号没有注册")
	ErrUserPwdError     = xerr.New(xerr.SERVER_COMMON_ERROR, "密码不正确")
)

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 基本参数校验：设备类型预处理
	in.DeviceType = strings.TrimSpace(in.DeviceType)
	if in.DeviceType == "" {
		return nil, errors.WithStack(xerr.New(xerr.REQUEST_PARAM_ERROR, "devicetype 不能为空"))
	}
	dt := strings.ToLower(in.DeviceType)
	if !constants.IsAllowedDevice(dt) {
		return nil, errors.WithStack(xerr.New(xerr.INVALID_DEVICE_TYPE, xerr.ErrMsg(xerr.INVALID_DEVICE_TYPE)))
	}
	in.DeviceType = dt

	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v , req %v", err, in.Phone)
	}

	// 密码验证
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire,
		userEntity.Id, in.DeviceType, in.DeviceName)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	}
	// 删除旧 token（幂等，不关心删除条数）
	if _, delErr := l.svcCtx.Redis.Del(userEntity.Id + ":" + in.DeviceType); delErr != nil {
		// 记录但不阻断登录，可考虑 metrics
		logx.WithContext(l.ctx).Errorf("redis del old token err: %v", delErr)
	}
	// 放入redis
	if setErr := l.svcCtx.Redis.Setex(userEntity.Id+":"+in.DeviceType, token, int(l.svcCtx.Config.Jwt.AccessExpire)); setErr != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "redis setex err %v", setErr)
	}
	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil

}
