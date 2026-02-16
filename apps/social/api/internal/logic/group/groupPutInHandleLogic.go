package group

import (
	"context"
	"github.com/pkg/errors"
	"slices"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/pkg/constants"
	"zeroIM/pkg/ctxdata"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupPutInHandleLogic 申请进群处理
func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleReq) (*types.GroupPutInHandleResp, error) {
	if req == nil || req.GroupReqId <= 0 || req.GroupId == "" {
		return nil, errors.WithStack(xerr.NewReqParamErr())
	}
	results := []constants.HandlerResult{constants.PassHandlerResult, constants.RejectHandlerResult, constants.CancelHandlerResult}
	if !slices.Contains(results, constants.HandlerResult(req.HandleResult)) {
		return nil, errors.WithStack(xerr.NewReqParamErr())
	}

	uid := ctxdata.GetUId(l.ctx)
	_, err := l.svcCtx.Social.GroupPutInHandle(l.ctx, &socialClient.GroupPutInHandleReq{
		GroupId:      req.GroupId,
		GroupReqId:   req.GroupReqId,
		HandleResult: req.HandleResult,
		HandleUid:    uid,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.GroupPutInHandleResp{}, nil
}
