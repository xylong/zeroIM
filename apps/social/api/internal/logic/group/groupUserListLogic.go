package group

import (
	"context"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"zeroIM/apps/social/api/internal/dto"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/apps/user/rpc/userClient"

	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupUserListLogic 成员列表列表
func NewGroupUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserListLogic {
	return &GroupUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserListLogic) GroupUserList(req *types.GroupUserListReq) (*types.GroupUserListResp, error) {
	// 获取群成员
	members, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialClient.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if members == nil || len(members.List) == 0 {
		return &types.GroupUserListResp{}, nil
	}

	// 获取用户信息
	users, err := l.svcCtx.User.FindUser(l.ctx, &userClient.FindUserReq{
		Ids: lo.Map(members.List, func(gm *socialClient.GroupMembers, _ int) string {
			return gm.UserId
		}),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userMap := make(map[string]*userClient.UserEntity, len(users.User))
	for _, user := range users.User {
		userMap[user.Id] = user
	}

	return &types.GroupUserListResp{
		List: dto.MemberToList(members, userMap),
	}, nil
}
