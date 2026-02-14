package logic

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutinListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinListLogic {
	return &GroupPutinListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutinListLogic) GroupPutinList(in *social.GroupPutinListReq) (*social.GroupPutinListResp, error) {
	reqs, err := l.svcCtx.Dao.GroupRequest.WithContext(l.ctx).
		Where(l.svcCtx.Dao.GroupRequest.GroupID.Eq(in.GroupId)).
		Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "list group req err %v groupId=%v", err, in.GroupId)
	}

	var list []*social.GroupRequests
	copier.Copy(&list, reqs)
	return &social.GroupPutinListResp{
		List: list,
	}, nil

}
