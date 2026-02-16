package dto

import (
	"zeroIM/apps/social/api/internal/types"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/apps/user/rpc/userClient"
)

func GroupToList(resp *socialClient.GroupListResp) []*types.Groups {
	if resp == nil || len(resp.List) == 0 {
		return nil
	}

	list := make([]*types.Groups, 0, len(resp.List))
	for _, v := range resp.List {
		list = append(list, &types.Groups{
			Id:              v.Id,
			Name:            v.Name,
			Icon:            v.Icon,
			Status:          int64(v.Status),
			GroupType:       int64(v.GroupType),
			IsVerify:        v.IsVerify,
			Notification:    v.Notification,
			NotificationUid: v.NotificationUid,
		})
	}

	return list
}

func GroupPutInToList(resp *socialClient.GroupPutinListResp) []*types.GroupRequests {
	if resp == nil || len(resp.List) == 0 {
		return nil
	}

	list := make([]*types.GroupRequests, 0, len(resp.List))
	for _, v := range resp.List {
		list = append(list, &types.GroupRequests{
			Id:            int64(v.Id),
			UserId:        v.ReqId,
			GroupId:       v.GroupId,
			ReqMsg:        v.ReqMsg,
			ReqTime:       v.ReqTime,
			JoinSource:    int64(v.JoinSource),
			InviterUserId: v.InviterUid,
			HandleResult:  int64(v.HandleResult),
			HandleUserId:  v.HandleUid,
			//HandleTime:    v.HandleTime,
		})
	}
	return list
}

func MemberToList(members *socialClient.GroupUsersResp, userMap map[string]*userClient.UserEntity) []*types.GroupMembers {
	list := make([]*types.GroupMembers, 0, len(members.List))

	for _, groupMembers := range members.List {
		member := types.GroupMembers{
			Id:          int64(groupMembers.Id),
			GroupId:     groupMembers.GroupId,
			UserId:      groupMembers.UserId,
			RoleLevel:   int(groupMembers.RoleLevel),
			InviterUid:  groupMembers.InviterUid,
			OperatorUid: groupMembers.OperatorUid,
		}

		// O(1) 查找用户信息
		if u, ok := userMap[member.UserId]; ok && u != nil {
			member.Nickname = u.Nickname
			member.UserAvatarUrl = u.Avatar
		}

		list = append(list, &member)
	}

	return list
}
