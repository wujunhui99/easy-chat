package middleware

import (
	"net/http"

	"github.com/wujunhui99/easy-chat/apps/social/api/internal/config"
	"github.com/wujunhui99/easy-chat/pkg/middleware/tokenmatch"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// JwtParseMiddleware 兼容层：保留原有引用点，内部直接委托给统一的 TokenMatch 中间件。
// 方便其它服务逐步替换时不需要立即修改所有 wiring。
type JwtParseMiddleware struct {
	delegate *tokenmatch.TokenMatch
}

func NewJwtParseMiddleware(c config.Config) *JwtParseMiddleware {
	tm := tokenmatch.New(redis.MustNewRedis(c.JwtTable), tokenmatch.Config{
		AccessSecret:    c.JwtAuth.AccessSecret,
		AllowEmptyToken: false,
	})
	return &JwtParseMiddleware{delegate: tm}
}

func (m *JwtParseMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return m.delegate.Handle(next)
}
