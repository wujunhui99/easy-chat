package logic

import (
	"context"
	"fmt"
	"time"

	rds "github.com/gomodule/redigo/redis"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/social/rpc/social"
	"github.com/wujunhui99/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type MutualFriendCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMutualFriendCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MutualFriendCountLogic {
	return &MutualFriendCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 共同好友数量（Redis 求交集）
func (l *MutualFriendCountLogic) MutualFriendCount(in *social.MutualFriendCountReq) (*social.MutualFriendCountResp, error) {
	uid := in.UserId
	oid := in.OtherId
	if uid == "" || oid == "" {
		return nil, xerr.New(xerr.REQUEST_PARAM_ERROR, xerr.ErrMsg(xerr.REQUEST_PARAM_ERROR))
	}
	if uid == oid {
		// 自己跟自己，共同好友为0
		return &social.MutualFriendCountResp{Count: 0}, nil
	}

	// Redis keys
	keyA := friendSetKey(uid)
	keyB := friendSetKey(oid)
	// 确保两个集合存在：若不存在则从 DB 构建（只放 status=0 的 friend_uid）
	if exists, _ := l.svcCtx.Redis.Exists(keyA); !exists {
		if err := l.buildFriendSetFromDB(uid, keyA); err != nil {
			return nil, err
		}
	}
	if exists, _ := l.svcCtx.Redis.Exists(keyB); !exists {
		if err := l.buildFriendSetFromDB(oid, keyB); err != nil {
			return nil, err
		}
	}

	// Redis 7.0+ SINTERCARD 直接返回交集基数
	// 优先尝试用 Eval 调用 SINTERCARD（需 Redis >= 7.0），失败则降级到 SINTERSTORE + SCARD
	if res, err := l.svcCtx.Redis.Eval(`return redis.call('SINTERCARD', 2, KEYS[1], KEYS[2])`, []string{keyA, keyB}); err == nil {
		if n, convErr := rds.Int64(res, nil); convErr == nil {
			return &social.MutualFriendCountResp{Count: n}, nil
		}
		return nil, fmt.Errorf("SINTERCARD: unexpected reply: %T (%v)", res, res)
	}

	// go-zero Redis 未直接封装 SINTERCARD，回退：SINTERSTORE 临时 key 再 SCARD
	tmpKey := fmt.Sprintf("tmp:mf:%s:%s:%d", uid, oid, time.Now().UnixNano())
	defer func() { _, _ = l.svcCtx.Redis.Del(tmpKey) }()

	if _, err := l.svcCtx.Redis.Sinterstore(tmpKey, keyA, keyB); err != nil {
		return nil, err
	}
	n, err := l.svcCtx.Redis.Scard(tmpKey)
	if err != nil {
		return nil, err
	}
	// 设置极短 ttl 便于快速回收
	_ = l.svcCtx.Redis.Expire(tmpKey, 3)

	return &social.MutualFriendCountResp{Count: int64(n)}, nil
}

func (l *MutualFriendCountLogic) buildFriendSetFromDB(uid, key string) error {
	// 读取该用户 status=0 的好友 friend_uid 列表
	// 利用已有 FriendsModel（字段含 user_id, friend_uid, status）
	// 这里没有现成的 ListByUserIdAndStatus 方法，我们直接用 QueryRowsNoCache 手写一条。
	// 为尽量少改动，直接写 SQL（也可考虑在 model 增加方法）。
	uids, err := l.svcCtx.FriendsModel.ListFriendUidsByUserId(l.ctx, uid)
	if err != nil {
		return err
	}
	if len(uids) == 0 {
		// 空集合：不写任何成员，交集计算自然为 0
		_ = l.svcCtx.Redis.Expire(key, 300)
		return nil
	}
	// 批量 SADD
	members := make([]interface{}, 0, len(uids))
	for _, fu := range uids {
		members = append(members, fu)
	}
	if _, err := l.svcCtx.Redis.Sadd(key, members...); err != nil {
		return err
	}
	// 设置 TTL
	_ = l.svcCtx.Redis.Expire(key, 300)
	return nil
}

func friendSetKey(uid string) string {
	return fmt.Sprintf("social:friends:set:%s", uid)
}
