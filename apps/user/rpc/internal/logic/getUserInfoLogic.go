package logic

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeroIM/apps/user/rpc/internal/svc"
	"zeroIM/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserNotExist = errors.New("用户不存在")
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	if l.svcCtx.Rdb != nil && in.Id != "" {
		key := "user:info:" + in.Id
		val, err := l.svcCtx.Rdb.Get(l.ctx, key).Result()
		if err == nil && len(val) > 0 {
			var cached user.UserEntity
			if json.Unmarshal([]byte(val), &cached) == nil {
				return &user.GetUserInfoResp{User: &cached}, nil
			}
		}
	}

	// 1.查数据
	userEntity, err := l.svcCtx.Dao.WithContext(l.ctx).User.
		Where(l.svcCtx.Dao.User.ID.Eq(in.Id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist
		}
		return nil, err
	}

	// 2.输出响应
	var resp user.UserEntity
	copier.Copy(&resp, userEntity)

	if l.svcCtx.Rdb != nil {
		key := "user:info:" + in.Id
		if b, err := json.Marshal(resp); err == nil {
			_ = l.svcCtx.Rdb.Set(l.ctx, key, b, 10*time.Minute).Err()
		}
	}

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
