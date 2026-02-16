package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zeroIM/apps/social/api/internal/logic/group"
	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"
)

// GroupPutInHandleHandler 申请进群处理
func GroupPutInHandleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupPutInHandleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewGroupPutInHandleLogic(r.Context(), svcCtx)
		resp, err := l.GroupPutInHandle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
