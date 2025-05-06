package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type RedisTokenCheckMiddleware struct {
	Redis *redis.Redis
}

func NewRedisTokenCheckMiddleware(rds *redis.Redis) *RedisTokenCheckMiddleware {
	return &RedisTokenCheckMiddleware{Redis: rds}
}

func (m *RedisTokenCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 取 Authorization header
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			httpx.Error(w, http.ErrNoCookie)
			return
		}
		var tokenString string 
		if strings.HasPrefix(authHeader, "Bearer "){
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}else{
			tokenString = authHeader
		}
		
		// 2. 解析 JWT，取 username
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			httpx.Error(w, err)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		username, _ := claims["username"].(string)
		if username == "" {
			httpx.Error(w, http.ErrNoCookie)
			return
		}

		// 3. 拼接两个 Redis key
		keyMobile := "user:" + username + ":mobile"
		keyDesktop := "user:" + username + ":desktop"

		ctx := r.Context()
		// 4. 分别查两个 key
		redisTokenMobile, _ := m.Redis.GetCtx(ctx, keyMobile)
		redisTokenDesktop, _ := m.Redis.GetCtx(ctx, keyDesktop)

		// 5. 对比
		if redisTokenMobile != tokenString && redisTokenDesktop != tokenString {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code": 401,
				"msg":  "登录已失效，请重新登录",
			})
			return
		}

		// 6. 校验通过，执行下一个 handler
		next(w, r)
	}
}
