package group

import (
	"net/http"

	"im-zero/east-chat/apps/social/api/internal/logic/group"
	"im-zero/east-chat/apps/social/api/internal/svc"
	"im-zero/east-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GroupPutInHandleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupPutInHandleRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := group.NewGroupPutInHandleLogic(r.Context(), svcCtx)
		resp, err := l.GroupPutInHandle(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}