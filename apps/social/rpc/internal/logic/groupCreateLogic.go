package logic

import (
	"context"
	"errors"
	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
	"zeroIM/apps/social/models"
	"zeroIM/apps/social/rpc/internal/dao"
	"zeroIM/pkg/wuid"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var (
	ErrGroupExists = xerr.NewMsgErr("群已存在")
)

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupCreate 建群
func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// 1.判断群是否存在
	group, err := l.svcCtx.Dao.Group.WithContext(l.ctx).
		Where(l.svcCtx.Dao.Group.CreatorUID.Eq(in.CreatorUid)).
		Where(l.svcCtx.Dao.Group.Name.Eq(in.Name)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "get group err %v req %v", err, in)
	}
	if group != nil {
		return nil, errors2.WithStack(ErrGroupExists)
	}

	// 2.建群
	l.svcCtx.Dao.Transaction(func(tx *dao.Query) error {
		var group = models.Group{
			ID:         wuid.GenUid(l.svcCtx.Config.Mysql.DSN),
			Name:       in.Name,
			Icon:       in.Icon,
			CreatorUID: in.CreatorUid,
			IsVerify:   0,
		}

		if err := tx.Group.WithContext(l.ctx).Create(&group); err != nil {
			return errors2.Wrapf(xerr.NewDBErr(), "create group err %v req %v", err, in)
		}

		now := time.Now()
		var groupMember = models.GroupMember{
			GroupID:    group.ID,
			UserID:     in.CreatorUid,
			RoleLevel:  1,
			JoinTime:   &now,
			JoinSource: 1,
		}
		if err := tx.GroupMember.WithContext(l.ctx).Create(&groupMember); err != nil {
			return errors2.Wrapf(xerr.NewDBErr(), "insert group member err %v req %v", err, in)
		}

		return nil
	})

	return &social.GroupCreateResp{}, nil
}
