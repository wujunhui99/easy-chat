/**
 * @author: wujunhui
 * @doc:
 */

package ctxdata

import (
	"context"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	DeveiceType contextKey = "devicetype"
	DeveiceName contextKey = "devicename"
	Identify    contextKey = "uid"
	// TokenRaw 在 API 层 JWT 解析后放入上下文，供后续 Redis 二次校验中间件复用，避免重复解析 Authorization
	TokenRaw contextKey = "tokenraw"
)

func GetJwtToken(secretKey string, iat, seconds int64, uid string, deveicetype string, devicename string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[string(Identify)] = uid
	claims[string(DeveiceType)] = deveicetype
	claims[string(DeveiceName)] = devicename

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}

// GetRawToken 读取已经标准化的原始 token 字符串
func GetRawToken(ctx context.Context) string {
	if v, ok := ctx.Value(TokenRaw).(string); ok {
		return v
	}
	return ""
}

func WithUid(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, Identify, uid)
}
func WithDeviceType(ctx context.Context, dev string) context.Context {
	return context.WithValue(ctx, DeveiceType, dev)
}
func WithRawToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, TokenRaw, token)
}
