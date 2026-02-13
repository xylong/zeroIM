package logic

import (
	"context"
	"errors"
	"time"
	"zeroIM/apps/social/models"
	"zeroIM/pkg/constants"
	"zeroIM/pkg/xerr"

	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"

	"zeroIM/apps/social/rpc/internal/dao"
	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
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
	if in.GroupId == "" || in.ReqId == "" {
		return nil, errors2.WithStack(xerr.NewReqParamErr())
	}

	// 2. 幂等性检查：是否已在群中
	member, err := l.svcCtx.Dao.GroupMember.WithContext(l.ctx).
		Where(l.svcCtx.Dao.GroupMember.GroupID.Eq(in.GroupId)).
		Where(l.svcCtx.Dao.GroupMember.UserID.Eq(in.ReqId)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find group member err %v, groupId %s, userId %s", err, in.GroupId, in.ReqId)
	}
	if member != nil {
		return &social.GroupPutinResp{GroupId: in.GroupId}, nil
	}

	// 3. 幂等性检查：是否已有未处理的申请
	req, err := l.svcCtx.Dao.GroupRequest.WithContext(l.ctx).
		Where(l.svcCtx.Dao.GroupRequest.GroupID.Eq(in.GroupId)).
		Where(l.svcCtx.Dao.GroupRequest.ReqID.Eq(in.ReqId)).
		Where(l.svcCtx.Dao.GroupRequest.HandleResult.Eq(int64(constants.NoHandlerResult))).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find group req err %v, groupId %s, userId %s", err, in.GroupId, in.ReqId)
	}
	if req != nil {
		return &social.GroupPutinResp{}, nil
	}

	// 4. 获取群信息
	group, err := l.svcCtx.Dao.Group.WithContext(l.ctx).
		Where(l.svcCtx.Dao.Group.ID.Eq(in.GroupId)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors2.Wrapf(xerr.NewDBErr(), "group not found, groupId %s", in.GroupId)
		}
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find group err %v, groupId %s", err, in.GroupId)
	}

	// 5. 准备申请记录
	groupReq := &models.GroupRequest{
		ReqID:         in.ReqId,
		GroupID:       in.GroupId,
		ReqMsg:        in.ReqMsg,
		ReqTime:       time.Now(),
		JoinSource:    int64(in.JoinSource),
		InviterUserID: in.InviterUid,
		HandleResult:  int64(constants.NoHandlerResult),
	}

	// 6. 判断是否可以自动通过
	isPass := false
	if group.GetVerify() {
		// 群开启了免验证进群
		isPass = true
	} else if constants.GroupJoinSource(in.JoinSource) == constants.InviteGroupJoinSource {
		// 邀请进群，检查邀请人权限
		inviter, err := l.svcCtx.Dao.GroupMember.WithContext(l.ctx).
			Where(l.svcCtx.Dao.GroupMember.UserID.Eq(in.InviterUid)).
			Where(l.svcCtx.Dao.GroupMember.GroupID.Eq(in.GroupId)).
			First()
		if err == nil {
			role := constants.GroupRoleLevel(inviter.RoleLevel)
			if role == constants.CreatorGroupRoleLevel || role == constants.ManagerGroupRoleLevel {
				isPass = true
				groupReq.HandleUserID = in.InviterUid
			}
		}
	}

	// 7. 执行入群逻辑（事务）
	if isPass {
		groupReq.HandleResult = int64(constants.PassHandlerResult)
		groupReq.HandleTime = func(t time.Time) *time.Time { return &t }(time.Now())

		err = l.svcCtx.Dao.Transaction(func(tx *dao.Query) error {
			if err := tx.GroupRequest.WithContext(l.ctx).Create(groupReq); err != nil {
				return err
			}
			return tx.GroupMember.WithContext(l.ctx).Create(&models.GroupMember{
				GroupID:     in.GroupId,
				UserID:      in.ReqId,
				OperatorUID: in.InviterUid,
				InviterUID:  in.InviterUid,
				RoleLevel:   int64(constants.AtLargeGroupRoleLevel),
				JoinTime:    time.Now(),
				JoinSource:  int64(in.JoinSource),
			})
		})
		if err != nil {
			return nil, errors2.Wrapf(xerr.NewDBErr(), "auto join group err %v, req %v", err, in)
		}
		return &social.GroupPutinResp{GroupId: in.GroupId}, nil
	}

	// 8. 需要验证，仅创建申请记录
	if err := l.svcCtx.Dao.GroupRequest.WithContext(l.ctx).Create(groupReq); err != nil {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "create group request err %v, req %v", err, in)
	}

	return &social.GroupPutinResp{}, nil
}
