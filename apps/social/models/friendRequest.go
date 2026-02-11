package models

import (
	"gorm.io/gorm"
	"time"
)

// FriendRequest 好友申请表
type FriendRequest struct {
	Id           int64          `gorm:"column:id;type:int(11) unsigned;primaryKey;autoIncrement:true" json:"id"`
	UserId       string         `gorm:"column:user_id;type:varchar(64);not null;comment:用户id" json:"user_id"`            // 用户id
	ReqUid       string         `gorm:"column:req_uid;type:varchar(64);not null;comment:申请好友id" json:"req_uid"`          // 申请好友id
	ReqMsg       string         `gorm:"column:req_msg;type:varchar(255);not null;comment:请求信息" json:"req_msg"`           // 请求信息
	ReqTime      time.Time      `gorm:"column:req_time;type:timestamp;not null;comment:请求时间" json:"req_time"`            // 请求时间
	HandleResult int64          `gorm:"column:handle_result;type:tinyint(4);not null;comment:处理结果" json:"handle_result"` // 处理结果
	HandleMsg    string         `gorm:"column:handle_msg;type:varchar(255);not null;comment:处理结果信息" json:"handle_msg"`   // 处理结果信息
	HandledAt    *time.Time     `gorm:"column:handled_at;type:timestamp;comment:处理时间" json:"handled_at"`                 // 处理时间
	CreatedAt    *time.Time     `gorm:"column:created_at;type:timestamp;not null;comment:创建时间" json:"created_at"`        // 创建时间
	UpdatedAt    *time.Time     `gorm:"column:updated_at;type:timestamp;not null;comment:更新时间" json:"updated_at"`        // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`                 // 删除时间
}

// TableName FriendRequest's table name
func (*FriendRequest) TableName() string {
	return "friend_requests"
}
