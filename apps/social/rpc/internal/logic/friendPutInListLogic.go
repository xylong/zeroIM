package logic

import (
	"context"
	"errors"
	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"
	"zeroIM/apps/social/models"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInListLogic) FriendPutInList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	reqs, err := l.svcCtx.Dao.FriendRequest.WithContext(l.ctx).
		Where(l.svcCtx.Dao.FriendRequest.UserId.Eq(in.UserId)).
		Find()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "list friend err %v req %v", err, in.UserId)
	}

	return &social.FriendPutInListResp{
		List: l.toList(reqs),
	}, nil
}

func (l *FriendPutInListLogic) toList(reqs []*models.FriendRequest) []*social.FriendRequests {
	if len(reqs) == 0 {
		return nil
	}

	var list []*social.FriendRequests
	for _, req := range reqs {
		list = append(list, &social.FriendRequests{
			Id:           int32(req.ID),
			UserId:       req.UserID,
			ReqUid:       req.ReqUID,
			ReqMsg:       req.ReqMsg,
			ReqTime:      req.CreatedAt.Unix(),
			HandleResult: int32(req.HandleResult),
		})
	}

	return list
}
