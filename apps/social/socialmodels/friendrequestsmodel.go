package socialmodels

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendRequestsModel = (*customFriendRequestsModel)(nil)

type (
	// FriendRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendRequestsModel.
	FriendRequestsModel interface {
		friendRequestsModel
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequests, error)
		ListNoHandler(ctx context.Context, userId string) ([]*FriendRequests, error)
		// ListByUserId 收到的(别人加我) => user_id = me 且未处理 (user_id=被申请人)
		ListByUserId(ctx context.Context, userId string) ([]*FriendRequests, error)
		// ListByReqUid 发出的(我加别人) => req_uid = me 且未处理 (req_uid=申请人)
		ListByReqUid(ctx context.Context, reqUid string) ([]*FriendRequests, error)
	}

	customFriendRequestsModel struct {
		*defaultFriendRequestsModel
	}
)

// NewFriendRequestsModel returns a model for the database table.
func NewFriendRequestsModel(conn sqlx.SqlConn, c cache.CacheConf) FriendRequestsModel {
	return &customFriendRequestsModel{
		defaultFriendRequestsModel: newFriendRequestsModel(conn, c),
	}
}

// Trans 事务封装（沿用之前逻辑）
func (m *customFriendRequestsModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

// FindByReqUidAndUserId 查询指定请求人与被请求人记录
func (m *customFriendRequestsModel) FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `user_id` = ?", friendRequestsRows, m.table)
	var resp FriendRequests
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, rid, uid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// ListNoHandler 列出待处理请求（handle_result=1）
func (m *customFriendRequestsModel) ListNoHandler(ctx context.Context, userId string) ([]*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `handle_result` = 1 and `user_id` = ?", friendRequestsRows, m.table)
	var resp []*FriendRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListByUserId 收到的(别人申请加我) user_id = me 且未处理
func (m *customFriendRequestsModel) ListByUserId(ctx context.Context, userId string) ([]*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `handle_result` = 1 and `user_id` = ? order by id desc", friendRequestsRows, m.table)
	var resp []*FriendRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListByReqUid 发出的(我申请加别人) req_uid = me 且未处理
func (m *customFriendRequestsModel) ListByReqUid(ctx context.Context, reqUid string) ([]*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `handle_result` = 1 and `req_uid` = ? order by id desc", friendRequestsRows, m.table)
	var resp []*FriendRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, reqUid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
