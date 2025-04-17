package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/junhui99/easy-chat/apps/im/ws/internal/svc"
	"github.com/junhui99/easy-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
)

type JwtAuth struct {
	svc    *svc.ServiceContext
	parser *token.TokenParser
	logx.Logger
}

func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	// fmt.Println("JwtAuth Auth")
	tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Errorf("parse token err %v ", err)
		return false
	}
	// fmt.Println("token")
	// fmt.Println("tok " ,tok.Raw)
	if !tok.Valid {
		fmt.Println("token not valid")
		return false
	}
	fmt.Println("token valid")
	// fmt.Printf("Claims type: %T, value: %+v\n", tok.Claims, tok.Claims)
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("token claims not ok")
		return false
	}

	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))

	return true

}

func (j *JwtAuth) UserId(r *http.Request) string {
	return ctxdata.GetUid(r.Context())
}
