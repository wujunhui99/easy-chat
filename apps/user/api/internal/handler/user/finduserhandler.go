package user

import (
	"net/http"

	userlogic "github.com/wujunhui99/easy-chat/apps/user/api/internal/logic/user"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 按手机号或ID查询用户(单条)
func FindUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFindReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := userlogic.NewFindUserLogic(r.Context(), svcCtx)
		resp, err := l.FindUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
