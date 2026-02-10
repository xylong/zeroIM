package user

import (
	"context"
	"zeroIM/apps/user/rpc/user"
	"zeroIM/pkg/ctxdata"

	"zeroIM/apps/user/api/internal/svc"
	"zeroIM/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDetailLogic 获取用户信息
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (*types.UserInfoResp, error) {
	uid := ctxdata.GetUId(l.ctx)

	userInfo, err := l.svcCtx.GetUserInfo(l.ctx, &user.GetUserInfoReq{Id: uid})
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResp{
		Info: types.User{
			Id:       userInfo.User.Id,
			Mobile:   userInfo.User.Phone,
			Nickname: userInfo.GetUser().Nickname,
			Sex:      byte(userInfo.User.Sex),
			Avatar:   userInfo.User.Avatar,
		},
	}, nil
}
