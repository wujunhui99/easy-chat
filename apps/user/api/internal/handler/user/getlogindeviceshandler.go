package user

import (
	"net/http"

	"github.com/wujunhui99/easy-chat/apps/user/api/internal/logic/user"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取当前账号已登录设备列表
func GetLoginDevicesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLoginDevicesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewGetLoginDevicesLogic(r.Context(), svcCtx)
		resp, err := l.GetLoginDevices(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
