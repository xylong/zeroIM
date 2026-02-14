package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"zeroIM/apps/user/models"
	"zeroIM/apps/user/rpc/internal/svc"
	"zeroIM/apps/user/rpc/user"
	"zeroIM/pkg/ctxdata"
	"zeroIM/pkg/encrypt"
	"zeroIM/pkg/logutil"
	"zeroIM/pkg/wuid"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

var (
	ErrPhoneIsRegister = errors.New("手机号已注册")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 1.验证手机号
	userModel := l.svcCtx.Dao.User
	_, err := userModel.WithContext(l.ctx).Where(userModel.Phone.Eq(in.Phone)).First()
	if err == nil {
		return nil, ErrPhoneIsRegister
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 2.入库
	nickname := in.Nickname
	if len(nickname) == 0 {
		// 默认昵称为手机号
		nickname = in.Phone
		// 或者隐藏中间4位: nickname = fmt.Sprintf("%s****%s", in.Phone[:3], in.Phone[7:])
		if len(in.Phone) >= 11 {
			nickname = fmt.Sprintf("%s****%s", in.Phone[:3], in.Phone[7:])
		}
	}

	avatar := in.Avatar
	if len(avatar) == 0 {
		// TODO: 设置默认头像URL
		avatar = ""
	}

	var (
		status = int8(1)
		sex    = int8(in.Sex)
	)

	userEntity := &models.User{
		ID:       wuid.GenUid(l.svcCtx.Config.Mysql.DSN),
		Avatar:   avatar,
		Nickname: nickname,
		Phone:    in.Phone,
		Status:   status,
		Sex:      sex,
	}
	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			logutil.BizError(l.ctx, "register.gen_password", err, map[string]interface{}{"phone": in.Phone})
			return nil, err
		}
		userEntity.Password = string(genPassword)
	}

	if err = userModel.WithContext(l.ctx).Create(userEntity); err != nil {
		logutil.BizError(l.ctx, "register.create_user", err, map[string]interface{}{"phone": in.Phone})
		return nil, err
	}

	// 3.生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.ID)
	if err != nil {
		logutil.BizError(l.ctx, "register.gen_token", err, map[string]interface{}{"uid": userEntity.ID})
		return nil, err
	}
	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
