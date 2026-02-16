package logic

import (
	"context"
	"errors"
	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"
	"zeroIM/apps/social/models"
	"zeroIM/apps/social/rpc/internal/dao"
	"zeroIM/pkg/constants"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendReqBeforePass   = xerr.NewMsgErr("好友申请已经通过")
	ErrFriendReqBeforeRefuse = xerr.NewMsgErr("好友申请已经被拒绝")
	ErrFriendReqCancelRefuse = xerr.NewMsgErr("好友申请已经被取消")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// 1.获取申请记录
	friendReq, err := l.svcCtx.Dao.FriendRequest.WithContext(l.ctx).Debug().
		Where(l.svcCtx.Dao.FriendRequest.Id.Eq(int64(in.FriendReqId))).
		Where(l.svcCtx.Dao.FriendRequest.UserId.Eq(in.UserId)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "get friendRequest by friendReqid err %v req %v", err, in.FriendReqId)
	}
	if friendReq == nil {
		return nil, errors2.WithStack(xerr.NewMsgErr("申请不存在"))
	}

	// 2.验证处理状态
	switch constants.HandlerResult(friendReq.HandleResult) {
	case constants.PassHandlerResult:
		return nil, errors2.WithStack(ErrFriendReqBeforePass)
	case constants.RejectHandlerResult:
		return nil, errors2.WithStack(ErrFriendReqBeforeRefuse)
	case constants.CancelHandlerResult:
		return nil, errors2.WithStack(ErrFriendReqCancelRefuse)
	}

	// 3.处理入库
	err = l.svcCtx.Dao.Transaction(func(tx *dao.Query) error {
		if _, err := tx.FriendRequest.WithContext(l.ctx).
			Where(tx.FriendRequest.Id.Eq(int64(in.FriendReqId))).
			Update(tx.FriendRequest.HandleResult, in.HandleResult); err != nil {
			return errors2.Wrapf(xerr.NewDBErr(), "update friendRequest err %v req %v", err, in)
		}

		// 如果通过。保存好友关系
		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		friends := []*models.Friend{
			{
				UserID:    friendReq.UserID,
				FriendUID: friendReq.ReqUID,
			}, {
				UserID:    friendReq.ReqUID,
				FriendUID: friendReq.UserID,
			},
		}
		if err := tx.Friend.WithContext(l.ctx).CreateInBatches(friends, 2); err != nil {
			return errors2.Wrapf(xerr.NewDBErr(), "create friend err %v req %v", err, in)
		}

		return nil
	})

	return &social.FriendPutInHandleResp{}, nil
}
