package logic

import (
	"context"
	"strconv"
	"time"
	"zeroIM/apps/user/models"
	"zeroIM/pkg/cachex"
	"zeroIM/pkg/xerr"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeroIM/apps/user/rpc/internal/svc"
	"zeroIM/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrInvalidId    = xerr.New(xerr.RequestParamError, "id错误")
	ErrUserNotExist = xerr.NewCodeErr(xerr.UserNotExist)

	cacheKeyPrefix = "user:info:"
	cacheTTL       = time.Minute * 60
	cacheNilTTL    = time.Second * 1 // 防止缓存穿透
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
	if err := l.validateRequest(in); err != nil {
		return nil, err
	}

	key := l.getCacheKey(in.Id)
	result, err := cachex.GetWithCache(l.ctx, l.svcCtx.Cache, key, cachex.Options{
		TTL:    cacheTTL,
		NilTTL: cacheNilTTL,
	}, func(ctx context.Context) (*user.GetUserInfoResp, error) {
		return l.getUserInfoInternal(in.Id)
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

// 参数校验
func (l *GetUserInfoLogic) validateRequest(in *user.GetUserInfoReq) error {
	if in == nil || in.Id <= 0 {
		return errors.WithStack(ErrInvalidId)
	}
	return nil
}

func (l *GetUserInfoLogic) getUserInfoInternal(uid int64) (*user.GetUserInfoResp, error) {
	userEntity, err := l.getUserFromDB(uid)
	if err != nil {
		return nil, err
	}

	return &user.GetUserInfoResp{
		User: l.toUserEntity(userEntity),
	}, nil
}

func (l *GetUserInfoLogic) getCacheKey(uid int64) string {
	return cacheKeyPrefix + strconv.Itoa(int(uid))
}

func (l *GetUserInfoLogic) getUserFromDB(uid int64) (*models.User, error) {
	userEntity, err := l.svcCtx.Dao.WithContext(l.ctx).User.Debug().
		Where(l.svcCtx.Dao.User.ID.Eq(uid)).
		First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist
		}

		logx.Errorf(
			"[GetUserInfo] db query failed, uid=%s, err=%v",
			uid, err,
		)
		return nil, err
	}

	return userEntity, nil
}

func (l *GetUserInfoLogic) toUserEntity(u *models.User) *user.UserEntity {
	return &user.UserEntity{
		Id:       u.ID,
		Avatar:   u.Avatar,
		Nickname: u.Nickname,
		Phone:    u.Phone,
		Status:   int32(u.Status),
		Sex:      int32(u.Sex),
	}
}
