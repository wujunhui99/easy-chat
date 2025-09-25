package tokenmatch

import (
	"context"
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
		ctx, authErr := t.Match(r)
		if authErr != nil {
			writeUnauthorized(w, authErr.Code, authErr.Msg)
			return
		}
		next(w, r.WithContext(ctx))
	}
}

// UnauthorizedError 表示鉴权阶段的业务错误，包含业务码和提示信息。
type UnauthorizedError struct {
	Code int
	Msg  string
}

func (e *UnauthorizedError) Error() string { return e.Msg }

func newUnauthorized(code int, msg string) *UnauthorizedError {
	return &UnauthorizedError{Code: code, Msg: msg}
}

func buildDeviceId(devType, devName string) string {
	if devName != "" {
		return devName
	}
	return devType
}

// Match 复用 JWT + Redis 一致性校验逻辑，供 HTTP 中间件和 WebSocket 鉴权复用。
// 成功时返回带有用户上下文的 ctx，失败时返回 UnauthorizedError。
func (t *TokenMatch) Match(r *http.Request) (context.Context, *UnauthorizedError) {
	ctx := r.Context()

	// 1. 读取 / 规范化 token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		if t.cfg.AllowEmptyToken { // 透传匿名
			return ctx, nil
		}
		return ctx, newUnauthorized(xerr.UNAUTHORIZED_ERROR, "缺少令牌")
	}
	tokenStr := authHeader
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		tokenStr = strings.TrimSpace(authHeader[7:])
	}
	if tokenStr == "" {
		return ctx, newUnauthorized(xerr.UNAUTHORIZED_ERROR, "令牌格式错误")
	}

	// 2. 解析 JWT（只做签名校验 + 基本字段提取）
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.cfg.AccessSecret), nil
	})
	if err != nil || !parsed.Valid {
		log.Println("tokenmatch: jwt parse invalid", err)
		return ctx, newUnauthorized(xerr.UNAUTHORIZED_ERROR, "令牌无效")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return ctx, newUnauthorized(xerr.UNAUTHORIZED_ERROR, "Claims 解析失败")
	}
	uidVal := claims[string(ctxdata.Identify)]
	devTypeVal := claims[string(ctxdata.DeveiceType)]
	devNameVal := claims[string(ctxdata.DeveiceName)]
	uid := strings.TrimSpace(fmt.Sprint(uidVal))
	devType := strings.TrimSpace(fmt.Sprint(devTypeVal))
	devName := strings.TrimSpace(fmt.Sprint(devNameVal))
	if uid == "" || devType == "" {
		return ctx, newUnauthorized(xerr.UNAUTHORIZED_ERROR, "缺少登录上下文")
	}
	deviceId := buildDeviceId(devType, devName)

	// 3. Redis 校验：确保 token 和当前注册的匹配
	redisToken, err := t.redis.Get(uid + ":" + devType)
	if err != nil {
		log.Println("tokenmatch: redis get err", err)
		return ctx, newUnauthorized(xerr.UNAUTHORIZED_ERROR, "状态验证失败")
	}
	if redisToken == "" { // 无记录：可能过期或被踢
		return ctx, newUnauthorized(xerr.UNAUTHORIZED_ERROR, "登录状态已过期")
	}
	if tokenStr != redisToken { // 不匹配：被其它登录覆盖
		return ctx, newUnauthorized(xerr.DEVICE_KICKED_ERROR, "当前设备登录状态已失效")
	}

	// 4. 注入上下文供后续逻辑使用
	ctx = ctxdata.WithUid(ctx, uid)
	ctx = ctxdata.WithDeviceType(ctx, devType)
	ctx = ctxdata.WithDeviceName(ctx, devName)
	ctx = ctxdata.WithDeviceId(ctx, deviceId)
	ctx = ctxdata.WithRawToken(ctx, tokenStr)
	fmt.Printf("uid is %s \n", uid)
	fmt.Printf("ctx uid is %s \n", ctx.Value("uid"))

	return ctx, nil
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
