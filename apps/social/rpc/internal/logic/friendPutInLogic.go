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

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// 1.是否已是好友
	friends, err := l.FindByUidAndFid(in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find friend err %v req %v", err, in)
	}
	if friends != nil {
		return &social.FriendPutInResp{}, errors2.Wrap(xerr.NewMsgErr("已是好友关系"), "")
	}

	// 2.是否有进行中或拒绝的申请
	request, err := l.FindByReqUidAndUserid(in.ReqUid, in.UserId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find friend request err %v req %v", err, in)
	}
	if request != nil {
		return &social.FriendPutInResp{}, errors2.Wrap(xerr.NewMsgErr("request already exist"), "")
	}
	// 3.入库
	err = l.svcCtx.Dao.FriendRequest.WithContext(l.ctx).Create(&models.FriendRequest{
		UserID: in.UserId,
		ReqUID: in.ReqUid,
		ReqMsg: in.ReqMsg,
	})
	if err != nil {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "create friendRequest err %v req %v", err, in)
	}

	return &social.FriendPutInResp{}, nil
}

func (l *FriendPutInLogic) FindByUidAndFid(uid, fid string) (*models.Friend, error) {
	result, err := l.svcCtx.Dao.Friend.WithContext(l.ctx).
		Where(l.svcCtx.Dao.Friend.UserId.Eq(uid)).
		Where(l.svcCtx.Dao.Friend.FriendUid.Eq(fid)).
		First()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (l *FriendPutInLogic) FindByReqUidAndUserid(reqUid, userId string) (*models.FriendRequest, error) {
	result, err := l.svcCtx.Dao.FriendRequest.WithContext(l.ctx).
		Where(l.svcCtx.Dao.FriendRequest.ReqUid.Eq(reqUid)).
		Where(l.svcCtx.Dao.FriendRequest.UserId.Eq(userId)).
		First()
	if err != nil {
		return nil, err
	}
	return result, nil
}
