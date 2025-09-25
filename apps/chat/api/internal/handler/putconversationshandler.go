package handler

import (
	"net/http"

	"github.com/wujunhui99/easy-chat/apps/chat/api/internal/logic"
	"github.com/wujunhui99/easy-chat/apps/chat/api/internal/svc"
	"github.com/wujunhui99/easy-chat/apps/chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func putConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewPutConversationsLogic(r.Context(), svcCtx)
		resp, err := l.PutConversations(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
