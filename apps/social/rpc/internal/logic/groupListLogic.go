package logic

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	groups, err := l.svcCtx.Dao.Group.WithContext(l.ctx).
		Where(l.svcCtx.Dao.Group.CreatorUID.Eq(in.UserId)).
		Find()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err %v req %v", err, in.UserId)
	}
	if len(groups) == 0 {
		return &social.GroupListResp{}, nil
	}

	var list []*social.Groups
	for _, group := range groups {
		list = append(list, &social.Groups{
			Id:              group.ID,
			Name:            group.Name,
			Icon:            group.Icon,
			Status:          int32(group.Status),
			CreatorUid:      group.CreatorUID,
			GroupType:       int32(group.GroupType),
			IsVerify:        group.GetVerify(),
			Notification:    group.Notification,
			NotificationUid: group.NotificationUID,
		})
	}
	return &social.GroupListResp{
		List: list,
	}, nil
}
