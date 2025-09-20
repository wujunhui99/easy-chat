package socialmodels

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendsModel = (*customFriendsModel)(nil)

type (
	// FriendsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendsModel.
	FriendsModel interface{
		friendsModel
		// ListFriendUidsByUserIdStatus0 列出某用户 status=0 的好友 friend_uid 列表
		ListFriendUidsByUserId(ctx context.Context, userId string) ([]string, error)
	}

	customFriendsModel struct {
		*defaultFriendsModel
	}
)

// NewFriendsModel returns a model for the database table.
func NewFriendsModel(conn sqlx.SqlConn, c cache.CacheConf) FriendsModel {
	return &customFriendsModel{
		defaultFriendsModel: newFriendsModel(conn, c),
	}
}

// ListFriendUidsByUserId 列出某用户 status=0 的好友 friend_uid 列表
func (m *customFriendsModel) ListFriendUidsByUserId(ctx context.Context, userId string) ([]string, error) {
	query := "select `friend_uid` from " + m.table + " where `user_id` = ? and `status` = 0"
	var rows []struct{
		FriendUid string `db:"friend_uid"`
	}
	if err := m.QueryRowsNoCacheCtx(ctx, &rows, query, userId); err != nil {
		return nil, err
	}
	uids := make([]string, 0, len(rows))
	for _, r := range rows {
		uids = append(uids, r.FriendUid)
	}
	return uids, nil
}
