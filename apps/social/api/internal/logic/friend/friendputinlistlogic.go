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

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	uid := ctxdata.GetUid(l.ctx)

	direction := req.Direction
	if direction != 2 { // 只认 2，其余归1
		direction = 1
	}

	list, err := l.svcCtx.Social.FriendPutInList(l.ctx, &socialclient.FriendPutInListReq{
		UserId:    uid,
		Direction: direction,
	})
	if err != nil {
		return nil, err
	}

	if len(list.List) == 0 {
		return &types.FriendPutInListResp{List: nil}, nil
	}

	ids := make([]string, 0, len(list.List)*2)
	seen := make(map[string]struct{}, len(list.List)*2)
	for _, item := range list.List {
		if item.UserId != "" {
			if _, ok := seen[item.UserId]; !ok {
				seen[item.UserId] = struct{}{}
				ids = append(ids, item.UserId)
			}
		}
		if item.ReqUid != "" {
			if _, ok := seen[item.ReqUid]; !ok {
				seen[item.ReqUid] = struct{}{}
				ids = append(ids, item.ReqUid)
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

	respList := make([]*types.FriendRequestView, 0, len(list.List))
	for _, item := range list.List {
		view := new(types.FriendRequestView)
		copier.Copy(view, item)
		if peer, ok := infoMap[item.UserId]; ok && peer != nil {
			view.UserNickname = peer.Nickname
			view.UserAvatar = peer.Avatar
		}
		if requester, ok := infoMap[item.ReqUid]; ok && requester != nil {
			view.ReqNickname = requester.Nickname
			view.ReqAvatar = requester.Avatar
		}
		respList = append(respList, view)
	}

	return &types.FriendPutInListResp{List: respList}, nil
}
