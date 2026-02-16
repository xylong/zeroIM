package group

import (
	"context"
	"github.com/pkg/errors"
	"zeroIM/apps/social/api/internal/dto"
	"zeroIM/apps/social/rpc/socialClient"

	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupPutInListLogic 申请进群列表
func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInListLogic) GroupPutInList(req *types.GroupPutInListReq) (*types.GroupPutInListResp, error) {
	resp, err := l.svcCtx.Social.GroupPutinList(l.ctx, &socialClient.GroupPutinListReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.GroupPutInListResp{
		List: dto.GroupPutInToList(resp),
	}, nil
}
