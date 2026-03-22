package constants

import (
	"encoding/json"
	"github.com/samber/lo"
)

// GroupType 群类型
type GroupType int8

const (
	GroupTypeNormal  GroupType = iota + 1 // 普通群
	GroupTypeCompany                      // 企业群
	GroupTypeChat                         // 聊天室
)

func (t GroupType) Valid() bool {
	switch t {
	case GroupTypeNormal, GroupTypeCompany, GroupTypeChat:
		return true
	}
	return false
}

func (t GroupType) String() string {
	switch t {
	case GroupTypeNormal:
		return "normal"
	case GroupTypeCompany:
		return "company"
	case GroupTypeChat:
		return "chat"
	default:
		return "unknown"
	}
}

func (t GroupType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// GroupStatus 群状态
type GroupStatus int8

const (
	GroupStatusOpen  GroupStatus = 1 // 开启
	GroupStatusClose GroupStatus = 0 // 关闭
)

func (l GroupStatus) Int8() int8 {
	return int8(l)
}

func (s GroupStatus) Valid() bool {
	switch s {
	case GroupStatusOpen, GroupStatusClose:
		return true
	}
	return false
}

// FriendHandlerResult 处理结果：0未处理 1通过 2拒绝 3取消
type FriendHandlerResult uint8

func (r FriendHandlerResult) Uint8() uint8 {
	return uint8(r)
}

const (
	FriendHandlerPending FriendHandlerResult = iota
	FriendHandlePass
	FriendHandleReject
	FriendHandleCancel
)

// GroupRoleLevel 群等级 0. 普通成员，10. 管理者，20. 群主
type GroupRoleLevel int8

func (l GroupRoleLevel) Uint8() uint8 {
	return uint8(l)
}

func (l GroupRoleLevel) IsAdmin() bool {
	return lo.Contains([]GroupRoleLevel{
		GroupOwner,
		GroupAdmin,
	}, l)
}

const (
	GroupOwner  GroupRoleLevel = 20
	GroupAdmin  GroupRoleLevel = 10
	GroupMember GroupRoleLevel = 0
)

// GroupJoinSource 进群申请的方式： 1. 邀请， 2. 申请
type GroupJoinSource int

func (l GroupJoinSource) Uint8() uint8 {
	return uint8(l)
}
func (l GroupJoinSource) Int32() int32 {
	return int32(l)
}

const (
	InviteGroupJoinSource GroupJoinSource = iota + 1
	PutInGroupJoinSource
)

// GroupVerify 群验证方式： 1. 开启， 0. 关闭
type GroupVerify int8

func (l GroupVerify) Int8() int8 {
	return int8(l)
}

const (
	GroupVerifyOpen  GroupVerify = 1
	GroupVerifyClose GroupVerify = 0
)

// GroupHandlerResult 入群申请处理结果：0待处理 1通过 2拒绝 3取消
type GroupHandlerResult uint8

func (l GroupHandlerResult) Uint8() uint8 {
	return uint8(l)
}

const (
	GroupHandlePending GroupHandlerResult = iota
	GroupHandlePass
	GroupHandleReject
	GroupHandleCancel
)
