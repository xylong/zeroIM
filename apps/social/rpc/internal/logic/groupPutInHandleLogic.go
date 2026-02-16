package logic

import (
	"context"
	"errors"
	"time"
	"zeroIM/apps/social/models"
	"zeroIM/apps/social/rpc/internal/dao"
	"zeroIM/pkg/constants"
	"zeroIM/pkg/xerr"

	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"

	"zeroIM/apps/social/rpc/internal/svc"
	"zeroIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrGroupReqBeforePass   = xerr.NewMsgErr("请求已通过")
	ErrGroupReqBeforeRefuse = xerr.NewMsgErr("请求已拒绝")
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	// 入群申请
	req, err := l.svcCtx.Dao.GroupRequest.WithContext(l.ctx).
		Where(l.svcCtx.Dao.GroupRequest.ID.Eq(int64(in.GroupReqId))).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors2.Wrap(xerr.NewMsgErr("group request not found"), "")
		}
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find group req err %v groupReqId=%v", err, in.GroupReqId)
	}

	// 申请审核人校验
	if in.GroupId != "" && req.GroupID != in.GroupId {
		return nil, errors2.Wrapf(xerr.NewReqParamErr(), "mismatch groupId %s != %s", in.GroupId, req.GroupID)
	}
	_, err = l.svcCtx.Dao.GroupMember.WithContext(l.ctx).
		Where(l.svcCtx.Dao.GroupMember.GroupID.Eq(in.GroupId)).
		Where(l.svcCtx.Dao.GroupMember.UserID.Eq(in.HandleUid)).
		Where(l.svcCtx.Dao.GroupMember.RoleLevel.In(int64(constants.CreatorGroupRoleLevel), int64(constants.ManagerGroupRoleLevel))).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors2.Wrapf(xerr.NewMsgErr("no permission"), "")
		}
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find group member err %v groupReqId=%v", err, in.GroupReqId)
	}

	switch constants.HandlerResult(req.HandleResult) {
	case constants.PassHandlerResult:
		return nil, errors2.WithStack(ErrGroupReqBeforePass)
	case constants.RejectHandlerResult:
		return nil, errors2.WithStack(ErrGroupReqBeforeRefuse)
	}
	req.HandleResult = int8(in.HandleResult)

	err = l.svcCtx.Dao.Transaction(func(tx *dao.Query) error {
		now := time.Now()
		if _, err := tx.GroupRequest.WithContext(l.ctx).
			Where(tx.GroupRequest.ID.Eq(int64(in.GroupReqId))).Updates(&models.GroupRequest{
			HandleUserID: in.HandleUid,
			HandleTime:   &now,
			HandleResult: int8(in.HandleResult),
		}); err != nil {
			return errors2.Wrapf(xerr.NewDBErr(), "update group req err %v groupReqId=%v", err, in.GroupReqId)
		}

		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		if err := tx.GroupMember.WithContext(l.ctx).Create(&models.GroupMember{
			GroupID:     req.GroupID,
			UserID:      req.ReqID,
			RoleLevel:   int8(constants.AtLargeGroupRoleLevel),
			OperatorUID: in.HandleUid,
		}); err != nil {
			return errors2.Wrapf(xerr.NewDBErr(), "create group member err %v groupReqId=%v", err, in.GroupReqId)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &social.GroupPutInHandleResp{}, nil
}
