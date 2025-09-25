package friend

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/api/internal/types"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/socialclient"
	"github.com/wujunhui99/easy-chat/apps/user/rpc/userclient"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	uid := ctxdata.GetUid(l.ctx)

	list, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{UserId: uid})
	if err != nil {
		return nil, err
	}

	if len(list.List) == 0 {
		return &types.FriendListResp{List: nil}, nil
	}

	ids := make([]string, 0, len(list.List))
	seen := make(map[string]struct{}, len(list.List))
	for _, item := range list.List {
		if item.FriendUid != "" {
			if _, ok := seen[item.FriendUid]; !ok {
				seen[item.FriendUid] = struct{}{}
				ids = append(ids, item.FriendUid)
			}
		}
	}

	infoMap := make(map[string]*userclient.UserEntity, len(ids))
	if len(ids) > 0 {
		if userResp, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{Ids: ids}); err != nil {
			l.Logger.Errorf("find user info failed for ids %v: %v", ids, err)
		} else if userResp != nil {
			for _, info := range userResp.User {
				if info == nil {
					continue
				}
				infoMap[info.Id] = info
			}
		}
	}

	respList := make([]*types.Friends, 0, len(list.List))
	for _, item := range list.List {
		friend := new(types.Friends)
		copier.Copy(friend, item)
		if info, ok := infoMap[item.FriendUid]; ok && info != nil {
			friend.Nickname = info.Nickname
			friend.Avatar = info.Avatar
		}
		respList = append(respList, friend)
	}

	return &types.FriendListResp{List: respList}, nil
}
