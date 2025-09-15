package user

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/types"
	userrpc "github.com/wujunhui99/easy-chat/apps/user/rpc/user"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前账号已登录设备列表
func NewGetLoginDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginDevicesLogic {
	return &GetLoginDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoginDevicesLogic) GetLoginDevices(req *types.GetLoginDevicesReq) (resp *types.GetLoginDevicesResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	if uid == "" {
		return nil, errors.WithStack(xerr.New(xerr.UNAUTHORIZED_ERROR, "缺少登录上下文"))
	}
	r, err := l.svcCtx.User.GetLoginDevices(l.ctx, &userrpc.GetLoginDevicesReq{CallerUid: uid})
	if err != nil {
		return nil, err
	}
	return &types.GetLoginDevicesResp{Devices: r.Devices}, nil
}
