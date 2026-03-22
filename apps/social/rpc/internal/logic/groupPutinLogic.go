package logic

import (
	"context"
	"errors"
	"fmt"
	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
	"zeroIM/apps/social/models"
	"zeroIM/apps/social/rpc/internal/dao"
	"zeroIM/pkg/cachex"
	"zeroIM/pkg/constants"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	GroupPutinCacheKey  = "group:putin:pending:%d:%d"
	GroupCacheKey       = "group:%d"
	GroupMemberCacheKey = "group:member:%d:%d"
)

var (
	GroupPutinExistErr  = xerr.NewMsgErr("重复申请")
	GroupMemberExistErr = xerr.NewMsgErr("已加入群")
)

type GroupPutinLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinLogic {
	return &GroupPutinLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {
	// 1. 参数校验
	if in.GroupId <= 0 || in.ReqId <= 0 || in.JoinSource <= 0 {
		return nil, errors2.WithStack(xerr.NewReqParamErr())
	}

	// 2. 幂等性检查：是否已在群中
	member, err := l.getGroupMember(in.GroupId, in.ReqId)
	if err != nil {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find group member err %v, groupId %d, userId %d", err, in.GroupId, in.ReqId)
	}
	if member != nil {
		return nil, GroupMemberExistErr
	}

	// 3. 幂等性检查：是否已申请
	req, err := l.getPendingGroupRequest(in.GroupId, in.ReqId)
	if err != nil {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find group request err %v, groupId %d, userId %d", err, in.GroupId, in.ReqId)
	}
	if req != nil {
		return nil, GroupPutinExistErr
	}

	// 4. 创建群申请
	group, err := l.getGroup(in.GroupId)
	if err != nil {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find group err %v, groupId %d", err, in.GroupId)
	}

	groupReq := &models.GroupRequest{
		ReqID:         in.ReqId,
		GroupID:       in.GroupId,
		ReqMsg:        in.ReqMsg,
		JoinSource:    uint8(in.JoinSource),
		InviterUserID: in.InviterUid,
	}

	isPass := false
	// 查询邀请人信息&权限
	if in.InviterUid > 0 {
		inviter, err := l.svcCtx.Dao.GroupMember.WithContext(l.ctx).
			Where(l.svcCtx.Dao.GroupMember.UserID.Eq(in.InviterUid)).
			Where(l.svcCtx.Dao.GroupMember.GroupID.Eq(in.GroupId)).
			First()
		if err != nil {
			return nil, errors2.Wrapf(xerr.NewDBErr(), "find group inviter err %v, groupId %d, userId %d", err, in.GroupId, in.InviterUid)
		}

		if constants.GroupRoleLevel(inviter.RoleLevel).IsAdmin() {
			isPass = true
			groupReq.HandleUserID = in.InviterUid
		}
	}

	if group.IsVerify == constants.GroupVerifyClose.Int8() {
		isPass = true
	}

	if isPass {
		groupReq.HandleResult = constants.GroupHandlePass.Uint8()
		groupReq.HandledAt = func(t time.Time) *time.Time { return &t }(time.Now())

		err = l.svcCtx.Dao.Transaction(func(tx *dao.Query) error {
			if err := tx.GroupRequest.WithContext(l.ctx).Create(groupReq); err != nil {
				return err
			}
			return tx.GroupMember.WithContext(l.ctx).Create(&models.GroupMember{
				GroupID:         in.GroupId,
				UserID:          in.ReqId,
				LastOperatorUID: in.InviterUid,
				InviterUID:      in.InviterUid,
				JoinSource:      uint8(in.JoinSource),
			})
		})
		if err != nil {
			return nil, errors2.Wrapf(xerr.NewDBErr(), "auto join group err %v, req %v", err, in)
		}
	} else {
		// 仅创建申请记录
		if err := l.svcCtx.Dao.GroupRequest.WithContext(l.ctx).Create(groupReq); err != nil {
			return nil, errors2.Wrapf(xerr.NewDBErr(), "create group request err %v, req %v", err, in)
		}
	}

	_ = l.svcCtx.Cache.Del(l.ctx, fmt.Sprintf(GroupPutinCacheKey, in.GroupId, in.ReqId))
	return &social.GroupPutinResp{GroupId: in.GroupId}, nil
}

func (l *GroupPutinLogic) getGroupMember(groupId, userId int64) (*models.GroupMember, error) {
	key := fmt.Sprintf(GroupMemberCacheKey, groupId, userId)

	return cachex.GetWithCache(
		l.ctx, l.svcCtx.Cache, key, cachex.Options{
			TTL:       time.Minute * 60,
			RandomTTL: time.Second * 100,
			CacheNil:  true,
		}, func(ctx context.Context) (*models.GroupMember, error) {
			member, err := l.svcCtx.Dao.GroupMember.WithContext(l.ctx).
				Where(l.svcCtx.Dao.GroupMember.GroupID.Eq(groupId)).
				Where(l.svcCtx.Dao.GroupMember.UserID.Eq(userId)).
				First()
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, nil
				}
				return nil, err
			}

			return member, nil
		})
}

func (l *GroupPutinLogic) getPendingGroupRequest(groupId, userId int64) (*models.GroupRequest, error) {
	key := fmt.Sprintf(GroupPutinCacheKey, groupId, userId)

	return cachex.GetWithCache(
		l.ctx, l.svcCtx.Cache, key, cachex.Options{
			TTL:       time.Minute * 60,
			RandomTTL: time.Second * 100,
			CacheNil:  true,
		}, func(ctx context.Context) (*models.GroupRequest, error) {
			req, err := l.svcCtx.Dao.GroupRequest.WithContext(l.ctx).Debug().
				Where(l.svcCtx.Dao.GroupRequest.GroupID.Eq(groupId)).
				Where(l.svcCtx.Dao.GroupRequest.ReqID.Eq(userId)).
				Where(l.svcCtx.Dao.GroupRequest.HandleResult.Eq(constants.GroupHandlePending.Uint8())).
				First()
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, nil
				}
				return nil, err
			}

			return req, nil
		})
}

func (l *GroupPutinLogic) getGroup(groupId int64) (*models.Group, error) {
	return cachex.GetWithCache(
		l.ctx, l.svcCtx.Cache, fmt.Sprintf(GroupCacheKey, groupId), cachex.Options{
			TTL:       time.Minute * 60,
			NilTTL:    time.Minute * 5,
			RandomTTL: time.Second * 100,
		}, func(ctx context.Context) (*models.Group, error) {
			group, err := l.svcCtx.Dao.Group.WithContext(l.ctx).Where(l.svcCtx.Dao.Group.ID.Eq(groupId)).First()
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, nil
				}
				return nil, err
			}
			return group, nil
		})
}
