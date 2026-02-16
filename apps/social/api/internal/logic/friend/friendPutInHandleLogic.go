package friend

import (
	"context"
	"fmt"
	errors2 "github.com/pkg/errors"
	"slices"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/pkg/constants"
	"zeroIM/pkg/ctxdata"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendPutInHandleLogic 好友申请处理
func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (*types.FriendPutInHandleResp, error) {
	if req == nil || req.FriendReqId <= 0 {
		fmt.Println(111)
		return nil, errors2.WithStack(xerr.NewReqParamErr())
	}
	results := []constants.HandlerResult{constants.PassHandlerResult, constants.RejectHandlerResult, constants.CancelHandlerResult}
	if !slices.Contains(results, constants.HandlerResult(req.HandleResult)) {
		return nil, errors2.WithStack(xerr.NewReqParamErr())
	}

	uid := ctxdata.GetUId(l.ctx)
	_, err := l.svcCtx.Social.FriendPutInHandle(l.ctx, &socialClient.FriendPutInHandleReq{
		FriendReqId:  req.FriendReqId,
		HandleResult: req.HandleResult,
		UserId:       uid,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("friend put in handle rpc failed: reqId=%s uid=%+v err=%+v", req.FriendReqId, uid, err)
		return nil, errors2.WithStack(err)
	}

	return &types.FriendPutInHandleResp{}, nil
}
