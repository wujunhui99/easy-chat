package user

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/types"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 退出登录
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, err error) {
	// todo: add your logic here and delete this line
	logoutResp, err := l.svcCtx.User.Logout(l.ctx, &user.LogoutReq{
		Id:         req.Id,
		DeviceType: req.DeviceType,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("login")
	var res types.LogoutResp
	copier.Copy(&res, logoutResp)

	return &res, nil
	return
}
