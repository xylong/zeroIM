package logic

import (
	"context"
	"fmt"
	"strings"
	"time"
	"zeroIM/pkg/constants"

	"errors"
	"zeroIM/apps/social/models"
	"zeroIM/pkg/cachex"
	"zeroIM/pkg/xerr"

	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	CacheFriendPrefix        = "social:friend:"
	CacheFriendRequestPrefix = "social:friend_req:"
	CacheTTL                 = time.Minute * 30
	CacheRandomTTL           = time.Second * 300
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
	if in.ReqUid == in.UserId {
		return nil, xerr.NewMsgErr("不能加自己为好友")
	}
	// 1.是否已有申请
	request, err := l.getFriendRequestWithCache(in.ReqUid, in.UserId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find friend request err %v req %v", err, in)
	}
	if request != nil {
		switch request.HandleResult {
		case constants.FriendHandlerPending.Uint8():
			return nil, xerr.NewMsgErr("申请已存在")
		case constants.FriendHandleReject.Uint8():
			return nil, xerr.NewMsgErr("申请已拒绝")
		case constants.FriendHandlePass.Uint8():
			return nil, xerr.NewMsgErr("申请已通过")
		}
		return &social.FriendPutInResp{}, nil
	}

	// 2.是否已是好友
	friends, err := l.getFriendWithCache(in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find friend err %v req %v", err, in)
	}
	if friends != nil {
		return nil, xerr.NewMsgErr("已是好友")
	}

	// 3.入库
	err = l.svcCtx.Dao.FriendRequest.WithContext(l.ctx).Create(&models.FriendRequest{
		UserID: in.UserId,
		ReqUID: in.ReqUid,
		ReqMsg: in.ReqMsg,
	})
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, xerr.NewMsgErr("申请已存在")
		}
		return nil, errors2.Wrapf(xerr.NewDBErr(), "create friendRequest err %v req %v", err, in)
	}
	// 4.删缓存
	_ = l.svcCtx.Cache.Del(l.ctx,
		l.getFriendRequestCacheKey(in.ReqUid, in.UserId),
		l.getFriendCacheKey(in.ReqUid, in.UserId),
	)

	return &social.FriendPutInResp{}, nil
}

func (l *FriendPutInLogic) getFriendWithCache(uid, fid int64) (*models.Friend, error) {
	key := l.getFriendCacheKey(uid, fid)

	query := func(ctx context.Context) (*models.Friend, error) {
		return l.FindByUidAndFid(uid, fid)
	}

	return cachex.GetWithCache(l.ctx, l.svcCtx.Cache, key, cachex.Options{
		TTL:       CacheTTL,
		RandomTTL: CacheRandomTTL,
	}, query)
}

func (l *FriendPutInLogic) getFriendRequestWithCache(reqUid, userId int64) (*models.FriendRequest, error) {
	key := l.getFriendRequestCacheKey(reqUid, userId)

	query := func(ctx context.Context) (*models.FriendRequest, error) {
		return l.FindByReqUidAndUserid(reqUid, userId)
	}

	return cachex.GetWithCache(l.ctx, l.svcCtx.Cache, key, cachex.Options{
		TTL:       CacheTTL,
		RandomTTL: CacheRandomTTL,
	}, query)
}

func (l *FriendPutInLogic) getFriendCacheKey(uid, fid int64) string {
	return fmt.Sprintf("%s%d:%d", CacheFriendPrefix, uid, fid)
}

func (l *FriendPutInLogic) getFriendRequestCacheKey(reqUid, userId int64) string {
	return fmt.Sprintf("%s%d:%d", CacheFriendRequestPrefix, reqUid, userId)
}

func (l *FriendPutInLogic) FindByUidAndFid(uid, fid int64) (*models.Friend, error) {
	result, err := l.svcCtx.Dao.Friend.WithContext(l.ctx).
		Where(l.svcCtx.Dao.Friend.UserID.Eq(uid)).
		Where(l.svcCtx.Dao.Friend.FriendUID.Eq(fid)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (l *FriendPutInLogic) FindByReqUidAndUserid(reqUid, userId int64) (*models.FriendRequest, error) {
	result, err := l.svcCtx.Dao.FriendRequest.WithContext(l.ctx).
		Where(l.svcCtx.Dao.FriendRequest.ReqUID.Eq(reqUid)).
		Where(l.svcCtx.Dao.FriendRequest.UserID.Eq(userId)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
