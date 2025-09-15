package user

import (
	"context"
	"errors"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/wujunhui99/easy-chat/apps/user/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/types"
	userrpc "github.com/wujunhui99/easy-chat/apps/user/rpc/user"
)

type FindUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// FindUser 按手机号或ID查询一个用户，若都为空返回错误；都传以手机号优先
func (l *FindUserLogic) FindUser(req *types.UserFindReq) (*types.UserFindResp, error) {
	phone := strings.TrimSpace(req.Phone)
	id := strings.TrimSpace(req.Id)
	if phone == "" && id == "" {
		return nil, errors.New("phone or id required")
	}

	rpcReq := &userrpc.FindUserReq{}
	if phone != "" {
		rpcReq.Phone = phone
	} else {
		rpcReq.Ids = []string{id}
	}

	r, err := l.svcCtx.User.FindUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	if len(r.User) == 0 {
		return &types.UserFindResp{User: nil}, nil
	}

	var u types.User
	copier.Copy(&u, r.User[0])
	return &types.UserFindResp{User: &u}, nil
}
