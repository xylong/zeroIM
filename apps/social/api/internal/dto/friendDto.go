package dto

import (
	"zeroIM/apps/social/api/internal/types"
	"zeroIM/apps/social/rpc/social"
	"zeroIM/apps/user/rpc/userClient"
)

func FriendToListResp(in *social.FriendListResp, records map[string]*userClient.UserEntity) []*types.Friends {
	if in == nil || len(in.List) == 0 {
		return nil
	}

	list := make([]*types.Friends, 0, len(in.List))
	for _, v := range in.List {
		friend := &types.Friends{
			Id:        v.Id,
			FriendUid: v.FriendUid,
		}

		if u, ok := records[v.FriendUid]; ok {
			friend.Nickname = u.Nickname
			friend.Avatar = u.Avatar
		}

		list = append(list, friend)
	}
	return list
}
