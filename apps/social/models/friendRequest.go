package models

import (
	"gorm.io/gorm"
	"time"
)

// FriendRequest 好友申请表
type FriendRequest struct {
	ID           uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       string         `gorm:"column:user_id;type:varchar(64);not null;index:idx_user_id_created_at,priority:1;comment:用户id" json:"user_id"`
	ReqUID       string         `gorm:"column:req_uid;type:varchar(64);not null;index:idx_req_uid_created_at,priority:1;comment:申请好友id" json:"req_uid"`
	ReqMsg       string         `gorm:"column:req_msg;type:varchar(255);not null;default:'';comment:请求信息" json:"req_msg"`
	HandleResult int8           `gorm:"column:handle_result;type:tinyint;not null;default:1;comment:处理结果：1-未处理 2-通过 3-拒绝 4-取消" json:"handle_result"`
	HandleMsg    string         `gorm:"column:handle_msg;type:varchar(255);not null;default:'';comment:处理结果信息" json:"handle_msg"`
	HandledAt    *time.Time     `gorm:"column:handled_at;type:timestamp;comment:处理时间" json:"handled_at,omitempty"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp;not null;index:idx_user_id_created_at,priority:2;index:idx_req_uid_created_at,priority:2" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName FriendRequest's table name
func (FriendRequest) TableName() string {
	return "friend_requests"
}
