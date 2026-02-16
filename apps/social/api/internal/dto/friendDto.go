package dto

import (
	"zeroIM/apps/social/api/internal/types"
	"zeroIM/apps/social/rpc/social"
	"zeroIM/apps/social/rpc/socialClient"
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

func FriendReqToListResp(in *socialClient.FriendPutInListResp) []*types.FriendRequests {
	if in == nil || len(in.List) == 0 {
		return nil
	}

	list := make([]*types.FriendRequests, 0, len(in.List))
	for _, v := range in.List {
		list = append(list, &types.FriendRequests{
			Id:           int64(v.Id),
			ReqUid:       v.ReqUid,
			ReqMsg:       v.ReqMsg,
			UserId:       v.UserId,
			HandleResult: int(v.HandleResult),
			ReqTime:      v.ReqTime,
		})
	}
	return list
}
