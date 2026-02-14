package logic

import (
	"context"
	"errors"
	errors2 "github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
	"zeroIM/pkg/ctxdata"
	"zeroIM/pkg/encrypt"
	"zeroIM/pkg/logutil"
	"zeroIM/pkg/xerr"

	"zeroIM/apps/user/rpc/internal/svc"
	"zeroIM/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister      = xerr.New(xerr.ServerCommonError, "手机号未注册")
	ErrUserPwdError          = xerr.New(xerr.ServerCommonError, "密码错误")
	ErrInvalidPasswordLength = xerr.New(xerr.RequestParamError, "密码长度错误")
	ErrInvalidPhone          = xerr.New(xerr.RequestParamError, "手机号格式错误")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	if len(in.Password) > 20 || len(in.Password) < 6 {
		return nil, ErrInvalidPasswordLength
	}
	if in.Phone == "" {
		return nil, ErrInvalidPhone
	}

	// 1.是否注册
	userEntity, err := l.svcCtx.Dao.WithContext(l.ctx).User.
		Where(l.svcCtx.Dao.User.Phone.Eq(in.Phone)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logutil.BizError(l.ctx, "login.phone_not_register", err, map[string]interface{}{"phone": in.Phone})
			return nil, errors2.WithStack(ErrPhoneNotRegister)
		}

		logutil.BizError(l.ctx, "login.query_user", err, map[string]interface{}{"phone": in.Phone})
		return nil, errors2.Wrapf(xerr.NewDBErr(), "find user by phone err %v, req %v", err, in.Phone)
	}

	// 2.密码校验
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password) {
		logutil.BizError(l.ctx, "login.user_password_error", err, map[string]interface{}{
			"id":       userEntity.ID,
			"phone":    in.Phone,
			"password": userEntity.Password,
		})
		return nil, errors2.WithStack(ErrUserPwdError)
	}

	// 3.生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.ID)
	if err != nil {
		return nil, errors2.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	}

	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
