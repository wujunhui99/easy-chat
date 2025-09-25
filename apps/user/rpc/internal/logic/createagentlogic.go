package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/wujunhui99/easy-chat/apps/user/models"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/user"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/wujunhui99/easy-chat/pkg/wuid"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrAgentPhoneExists = errors.New("agent phone already exists")

type CreateAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentLogic {
	return &CreateAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentLogic) CreateAgent(in *user.CreateAgentReq) (*user.CreateAgentResp, error) {
	userId := wuid.GenUid(l.svcCtx.Config.Mysql.DataSource)
	nickname := in.GetNickname()
	if nickname == "" {
		nickname = fmt.Sprintf("Agent-%s", userId[len(userId)-4:])
	}
	phone := in.GetPhone()
	if phone == "" {
		phone = fmt.Sprintf("agent_%s", userId)
	} else {
		exists, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, phone)
		if err != nil && err != models.ErrNotFound {
			return nil, err
		}
		if exists != nil {
			return nil, ErrAgentPhoneExists
		}
	}

	userEntity := &models.Users{
		Id:       userId,
		Avatar:   in.GetAvatar(),
		Nickname: nickname,
		Phone:    phone,
		Status: sql.NullInt64{
			Int64: constants.UserStatusNormal,
			Valid: true,
		},
		UserType: constants.UserTypeAgent,
	}

	if in.GetSex() != 0 {
		userEntity.Sex = sql.NullInt64{
			Int64: int64(in.GetSex()),
			Valid: true,
		}
	}

	if _, err := l.svcCtx.UsersModel.Insert(l.ctx, userEntity); err != nil {
		return nil, err
	}

	return &user.CreateAgentResp{Id: userId}, nil
}
