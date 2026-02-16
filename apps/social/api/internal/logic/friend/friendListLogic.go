package friend

import (
	"context"
	"github.com/pkg/errors"
	"zeroIM/apps/social/api/internal/dto"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/apps/user/rpc/userClient"
	"zeroIM/pkg/ctxdata"

	"zeroIM/apps/social/api/internal/svc"
	"zeroIM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendListLogic 好友列表
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (*types.FriendListResp, error) {
	// 获取好友
	friends, err := l.svcCtx.Social.FriendList(l.ctx, &socialClient.FriendListReq{UserId: ctxdata.GetUId(l.ctx)})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if friends == nil || len(friends.List) == 0 {
		return &types.FriendListResp{}, nil
	}

	// 获取好友信息
	uids := make([]string, 0, len(friends.List))
	for _, f := range friends.List {
		uids = append(uids, f.FriendUid)
	}
	users, err := l.svcCtx.User.FindUser(l.ctx, &userClient.FindUserReq{
		Ids: uids,
	})

	if err != nil {
		return &types.FriendListResp{}, errors.WithStack(err)
	}
	userRecords := make(map[string]*userClient.UserEntity, len(users.User))
	for i, _ := range users.User {
		userRecords[users.User[i].Id] = users.User[i]
	}

	return &types.FriendListResp{
		List: dto.FriendToListResp(friends, userRecords),
	}, nil
}
