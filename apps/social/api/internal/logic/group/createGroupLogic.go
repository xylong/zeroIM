package group

import (
	"context"
	"github.com/pkg/errors"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/pkg/ctxdata"

	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateGroupLogic 创群
func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (*types.GroupCreateResp, error) {
	uid := ctxdata.GetUId(l.ctx)
	_, err := l.svcCtx.Social.GroupCreate(l.ctx, &socialClient.GroupCreateReq{
		CreatorUid: uid,
		Name:       req.Name,
		Icon:       req.Icon,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &types.GroupCreateResp{}, nil
}
