package logic

import (
	"context"
	"github.com/jinzhu/copier"
	errors2 "github.com/pkg/errors"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	friends, err := l.svcCtx.Dao.Friend.WithContext(l.ctx).
		Where(l.svcCtx.Dao.Friend.UserId.Eq(in.UserId)).
		Find()
	if err != nil {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "list friend by uid err %v uid %v", err, in.UserId)
	}

	var listResp []*social.Friends
	copier.Copy(&listResp, friends)

	return &social.FriendListResp{
		List: listResp,
	}, nil
}
