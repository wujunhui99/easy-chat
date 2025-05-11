package middleware

import (
	"log"
	"net/http"

	"github.com/wujunhui99/easy-chat/apps/user/api/internal/config"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type contextKey string

type JwtParseMiddleware struct {
	secret string
	*redis.Redis
}

func NewJwtParseMiddleware(c config.Config) *JwtParseMiddleware {
	return &JwtParseMiddleware{secret: c.JwtAuth.AccessSecret,
		Redis: redis.MustNewRedis(c.JwtTable)}
}

func (m *JwtParseMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		uid := ctxdata.GetUid(r.Context())
		devicetype := ctxdata.GetDevicetype(r.Context())
		rediskey := uid + ":" + devicetype
		redistoken, err := m.Redis.Get(rediskey)
		if err != nil {
			log.Println("redis get token err ", err)
			return
		}
		token := r.Header.Get("Authorization")
		if token != redistoken {
			log.Println("token not match")
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		// Passthrough to next handler if need
		next(w, r)
	}
}
