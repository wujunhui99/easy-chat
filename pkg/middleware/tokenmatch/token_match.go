package tokenmatch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/wujunhui99/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// Config 公共中间件配置
type Config struct {
	// JWT 秘钥
	AccessSecret string
	// 允许无 token 透传（少量公开接口可以复用同一链）
	AllowEmptyToken bool
}

// TokenMatch 统一：提取 Authorization -> 解析 JWT -> 注入上下文(uid/deviceType/rawToken) -> Redis 比对 token 一致性
// 目标：跨多个 API 服务直接复用；避免再写独立 JwtParseMiddleware。
type TokenMatch struct {
	redis *redis.Redis
	cfg   Config
}

func New(r *redis.Redis, cfg Config) *TokenMatch { return &TokenMatch{redis: r, cfg: cfg} }

// Handle : 标准 go-zero middleware
func (t *TokenMatch) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 读取 / 规范化 token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			if t.cfg.AllowEmptyToken { // 透传匿名
				next(w, r)
				return
			}
			writeUnauthorized(w, xerr.UNAUTHORIZED_ERROR, "缺少令牌")
			return
		}
		tokenStr := authHeader
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			tokenStr = strings.TrimSpace(authHeader[7:])
		}
		if tokenStr == "" {
			writeUnauthorized(w, xerr.UNAUTHORIZED_ERROR, "令牌格式错误")
			return
		}

		// 2. 解析 JWT（只做签名校验 + 基本字段提取）
		parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) { return []byte(t.cfg.AccessSecret), nil })
		if err != nil || !parsed.Valid {
			log.Println("tokenmatch: jwt parse invalid", err)
			writeUnauthorized(w, xerr.UNAUTHORIZED_ERROR, "令牌无效")
			return
		}
		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			writeUnauthorized(w, xerr.UNAUTHORIZED_ERROR, "Claims 解析失败")
			return
		}
		uid, _ := claims[string(ctxdata.Identify)].(string)
		devType, _ := claims[string(ctxdata.DeveiceType)].(string)
		if uid == "" || devType == "" {
			writeUnauthorized(w, xerr.UNAUTHORIZED_ERROR, "0c缺少登录上下文")
			return
		}

		// 3. Redis 校验：确保 token 和当前注册的匹配
		redisToken, err := t.redis.Get(uid + ":" + devType)
		if err != nil {
			log.Println("tokenmatch: redis get err", err)
			writeUnauthorized(w, xerr.UNAUTHORIZED_ERROR, "状态验证失败")
			return
		}
		if redisToken == "" { // 无记录：可能过期或被踢
			writeUnauthorized(w, xerr.UNAUTHORIZED_ERROR, "登录状态已过期")
			return
		}
		if tokenStr != redisToken { // 不匹配：被其它登录覆盖
			writeUnauthorized(w, xerr.DEVICE_KICKED_ERROR, "当前设备登录状态已失效")
			return
		}

		// 4. 注入上下文供后续逻辑使用
		ctx := r.Context()
		ctx = ctxdata.WithUid(ctx, uid)
		ctx = ctxdata.WithDeviceType(ctx, devType)
		ctx = ctxdata.WithRawToken(ctx, tokenStr)
		r = r.WithContext(ctx)
		fmt.Printf("uid is %s \n", uid)
		fmt.Printf("ctx uid is %s \n", ctx.Value("uid"))

		next(w, r)
	}
}

// writeUnauthorized 统一未授权响应
func writeUnauthorized(w http.ResponseWriter, bizCode int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code": bizCode,
		"msg":  msg,
		"data": nil,
	})
}
