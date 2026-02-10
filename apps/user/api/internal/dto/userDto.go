package dto

import (
	"zeroIM/apps/user/api/internal/types"
	"zeroIM/apps/user/rpc/user"
)

func UserInfoToResp(in *user.GetUserInfoResp) *types.UserInfoResp {
	if in == nil || in.User == nil {
		return nil
	}

	return &types.UserInfoResp{
		Info: types.User{
			Id:       in.User.Id,
			Mobile:   in.User.Phone,
			Nickname: in.User.Nickname,
			Sex:      byte(in.User.Sex),
			Avatar:   in.User.Avatar,
		},
	}
}
