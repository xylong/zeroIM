package logic

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"
	"zeroIM/apps/social/models"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUsersLogic {
	return &GroupUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupUsersLogic) GroupUsers(in *social.GroupUsersReq) (*social.GroupUsersResp, error) {
	members, err := l.svcCtx.Dao.GroupMember.WithContext(l.ctx).
		Where(l.svcCtx.Dao.GroupMember.GroupID.Eq(in.GroupId)).
		Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "group member by groupid err %v groupid %v", err, in.GroupId)
	}

	return &social.GroupUsersResp{
		List: l.toList(members),
	}, nil
}

func (l *GroupUsersLogic) toList(members []*models.GroupMember) []*social.GroupMembers {
	if len(members) == 0 {
		return nil
	}

	var list []*social.GroupMembers
	_ = copier.Copy(&list, members)
	return list
}
