package user

import (
	"fmt"
	"net/http"

	"github.com/wujunhui99/easy-chat/apps/user/api/internal/logic/user"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户登入
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("LoginHandler....")
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
