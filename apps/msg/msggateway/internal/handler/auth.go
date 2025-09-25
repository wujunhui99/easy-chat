package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wujunhui99/easy-chat/apps/msg/msggateway/internal/svc"
	"github.com/wujunhui99/easy-chat/pkg/ctxdata"
	"github.com/wujunhui99/easy-chat/pkg/middleware/tokenmatch"
	"github.com/zeromicro/go-zero/core/logx"
)

type JwtAuth struct {
	svc *svc.ServiceContext
	logx.Logger
}

func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		Logger: logx.WithContext(context.Background()),
	}
}

func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	ctx, err := j.svc.TokenMatch.Match(r)
	if err != nil {
		j.Errorf("token match failed code=%d msg=%s", err.Code, err.Msg)
		writeWsUnauthorized(w, err)
		return false
	}
	*r = *r.WithContext(ctx)
	return true
}

func (j *JwtAuth) UserId(r *http.Request) string {
	return ctxdata.GetUid(r.Context())
}

func (j *JwtAuth) DeviceId(r *http.Request) string {
	if id := ctxdata.GetDeviceId(r.Context()); id != "" {
		return id
	}
	return ctxdata.GetDevicetype(r.Context())
}

func writeWsUnauthorized(w http.ResponseWriter, err *tokenmatch.UnauthorizedError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code": err.Code,
		"msg":  err.Msg,
		"data": nil,
	})
}
