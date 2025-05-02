package logic

import (
	"context"
	"time"

	"github.com/junhui99/easy-chat/apps/user/models"
	"github.com/junhui99/easy-chat/apps/user/rpc/internal/svc"
	"github.com/junhui99/easy-chat/apps/user/rpc/user"
	"github.com/junhui99/easy-chat/pkg/ctxdata"
	"github.com/junhui99/easy-chat/pkg/encrypt"
	"github.com/junhui99/easy-chat/pkg/xerr"
	"github.com/pkg/errors"

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
	// todo: add your logic here and delete this line
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
		userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	}
	// 删除之外的token
	_, err = l.svcCtx.Redis.Del(userEntity.Id + ":" + in.DeviceType)
	// 放入redis
	err = l.svcCtx.Redis.Setex(userEntity.Id+":"+in.DeviceType, token, int(l.svcCtx.Config.Jwt.AccessExpire))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "redis setex err %v", err)
	}
	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil

}
