package group

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupPutInLogic 申请进群
func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInLogic) GroupPutIn(req *types.GroupPutInReq) (*types.GroupPutInResp, error) {
	uid := ctxdata.GetUId(l.ctx)
	_, err := l.svcCtx.Social.GroupPutin(l.ctx, &socialClient.GroupPutinReq{
		GroupId:    req.GroupId,
		ReqId:      uid,
		JoinSource: int32(req.JoinSource),
		ReqMsg:     req.ReqMsg,
	})
	if err != nil {
		fmt.Println("RPC ERROR:", err)
		return nil, errors.WithStack(err)
	}

	return &types.GroupPutInResp{}, nil
}
