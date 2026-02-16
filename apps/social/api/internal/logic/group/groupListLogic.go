package group

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

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupListLogic 用户申群列表
func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupListLogic) GroupList(req *types.GroupListReq) (*types.GroupListResp, error) {
	uid := ctxdata.GetUId(l.ctx)
	resp, err := l.svcCtx.Social.GroupList(l.ctx, &socialClient.GroupListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.GroupListResp{
		List: dto.GroupToList(resp),
	}, nil
}
