package dto

import (
	"zeroIM/apps/social/api/internal/types"
	"zeroIM/apps/social/rpc/socialClient"
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
