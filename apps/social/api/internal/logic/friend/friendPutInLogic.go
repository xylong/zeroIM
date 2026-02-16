package friend

import (
	"context"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/pkg/ctxdata"
	"zeroIM/pkg/xerr"

	errors2 "github.com/pkg/errors"

	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendPutInLogic 好友申请
func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (*types.FriendPutInResp, error) {
	uid := ctxdata.GetUId(l.ctx)
	if req == nil || req.UserId == "" {
		return nil, errors2.WithStack(xerr.NewReqParamErr())
	}
	if uid == "" {
		return nil, errors2.WithStack(xerr.NewCodeErr(xerr.TokenExpireError))
	}
	if req.UserId == uid {
		return nil, errors2.WithStack(xerr.NewMsgErr("不能添加自己为好友"))
	}
	if len(req.ReqMsg) > 256 {
		return nil, errors2.WithStack(xerr.NewReqParamErr())
	}

	rpcReq := &socialClient.FriendPutInReq{
		UserId: uid,
		ReqUid: req.UserId,
		ReqMsg: req.ReqMsg,
	}
	_, err := l.svcCtx.Social.FriendPutIn(l.ctx, rpcReq)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("friend put in rpc failed: uid=%s req=%+v err=%+v", uid, req, err)
		return nil, errors2.WithStack(err)
	}

	return &types.FriendPutInResp{}, nil
}
