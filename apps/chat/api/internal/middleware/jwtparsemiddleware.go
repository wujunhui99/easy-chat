package middleware

import (
	"net/http"

	"github.com/wujunhui99/easy-chat/apps/chat/api/internal/config"
	"github.com/wujunhui99/easy-chat/pkg/middleware/tokenmatch"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type JwtParseMiddleware struct{ tm *tokenmatch.TokenMatch }

func NewJwtParseMiddleware(c config.Config) *JwtParseMiddleware {
	return &JwtParseMiddleware{tm: tokenmatch.New(
		redis.MustNewRedis(c.JwtTable),
		tokenmatch.Config{
			AccessSecret:    c.JwtAuth.AccessSecret,
			AllowEmptyToken: false,
		},
	)}
}

func (m *JwtParseMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc { return m.tm.Handle(next) }
