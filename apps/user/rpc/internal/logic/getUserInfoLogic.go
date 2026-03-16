package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"zeroIM/apps/user/models"
	"zeroIM/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"zeroIM/apps/user/rpc/internal/svc"
	"zeroIM/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrInvalidId    = xerr.New(xerr.RequestParamError, "id错误")
	ErrUserNotExist = xerr.NewCodeErr(xerr.UserNotExist)

	cacheKeyPrefix = "user:info:"
	cacheTTL       = time.Minute * 10
	cacheNilTTL    = time.Second * 15 // 防止缓存穿透
	cacheNilValue  = "__nil__"
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
	// 1.参数校验
	if err := l.validateRequest(in); err != nil {
		return nil, err
	}
	// 2. 使用 SingleFlight 防止缓存击穿
	key := l.getCacheKey(in.Id)
	result, err := l.svcCtx.UserInfoSF.Do(key, func() (interface{}, error) {
		return l.getUserInfoInternal(in.Id)
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result.(*user.GetUserInfoResp), nil
}

// 参数校验
func (l *GetUserInfoLogic) validateRequest(in *user.GetUserInfoReq) error {
	if in == nil || in.Id <= 0 {
		return errors.WithStack(ErrInvalidId)
	}
	return nil
}

func (l *GetUserInfoLogic) getUserInfoInternal(uid int64) (*user.GetUserInfoResp, error) {
	// 1.从缓存获取
	if u, hit, err := l.getUserFromCache(uid); err == nil && hit {
		if u == nil {
			return nil, ErrUserNotExist
		}
		return &user.GetUserInfoResp{User: u}, nil
	}

	// 2.从db查
	userEntity, err := l.getUserFromDB(uid)
	if err != nil {
		return nil, err
	}

	// 3.转换响应
	resp := l.toUserEntity(userEntity)

	// 4.异步写缓存
	l.setUserCacheAsync(uid, resp)

	return &user.GetUserInfoResp{
		User: resp,
	}, nil
}

func (l *GetUserInfoLogic) getCacheKey(uid int64) string {
	return cacheKeyPrefix + strconv.Itoa(int(uid))
}

func (l *GetUserInfoLogic) getUserFromCache(uid int64) (*user.UserEntity, bool, error) {
	if l.svcCtx.Rdb == nil {
		return nil, false, nil
	}

	key := l.getCacheKey(uid)
	val, err := l.svcCtx.Rdb.Get(l.ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		logx.Errorf("[GetUserInfo] redis get failed, key=%s, err=%v", key, err)
		return nil, false, err
	}

	if val == cacheNilValue {
		return nil, true, nil
	}

	var u user.UserEntity
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		logx.Errorf("[GetUserInfo] redis unmarshal failed, key=%s, err=%v", key, err)
		return nil, false, err
	}

	return &u, true, nil
}

func (l *GetUserInfoLogic) getUserFromDB(uid int64) (*models.User, error) {
	uidStr := strconv.Itoa(int(uid))
	userEntity, err := l.svcCtx.Dao.WithContext(l.ctx).User.Debug().
		Where(l.svcCtx.Dao.User.ID.Eq(uidStr)).
		First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.cacheNilValue(uidStr)
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

func (l *GetUserInfoLogic) setUserCacheAsync(uid int64, u *user.UserEntity) {
	if l.svcCtx.Rdb == nil {
		return
	}

	go func() {
		b, _ := json.Marshal(u)
		_ = l.svcCtx.Rdb.Set(
			context.Background(),
			cacheKeyPrefix+strconv.Itoa(int(uid)),
			b,
			cacheTTL,
		).Err()
	}()
}

func (l *GetUserInfoLogic) cacheNilValue(uid string) {
	if l.svcCtx.Rdb == nil {
		return
	}

	go func() {
		_ = l.svcCtx.Rdb.Set(
			context.Background(),
			cacheKeyPrefix+uid,
			cacheNilValue,
			cacheNilTTL,
		).Err()
	}()
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
