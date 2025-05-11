package user

import (
	"net/http"

	"github.com/wujunhui99/easy-chat/apps/user/api/internal/logic/user"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 退出登录
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LogoutReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewLogoutLogic(r.Context(), svcCtx)
		resp, err := l.Logout(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
