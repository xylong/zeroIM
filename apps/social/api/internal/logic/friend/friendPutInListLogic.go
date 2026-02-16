package friend

import (
	"context"
	"github.com/pkg/errors"
	"zeroIM/apps/social/api/internal/dto"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/pkg/ctxdata"

	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendPutInListLogic 好友申请列表
func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (*types.FriendPutInListResp, error) {
	requests, err := l.svcCtx.Social.FriendPutInList(l.ctx, &socialClient.FriendPutInListReq{
		UserId: ctxdata.GetUId(l.ctx),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.FriendPutInListResp{
		List: dto.FriendReqToListResp(requests),
	}, nil
}
