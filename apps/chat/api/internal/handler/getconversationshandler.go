package handler

import (
	"net/http"

	"github.com/wujunhui99/easy-chat/apps/chat/api/internal/logic"
	"github.com/wujunhui99/easy-chat/apps/chat/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetConversationsLogic(r.Context(), svcCtx)
		resp, err := l.GetConversations(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
