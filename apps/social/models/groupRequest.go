package models

import (
	"gorm.io/gorm"
	"time"
)

// GroupRequest 入群申请表
type GroupRequest struct {
	ID            uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ReqID         string         `gorm:"column:req_id;type:varchar(64);not null;index:idx_req_created,priority:1;comment:申请人uid" json:"req_id"`
	GroupID       string         `gorm:"column:group_id;type:varchar(64);not null;index:idx_group_created,priority:1;comment:群id" json:"group_id"`
	ReqMsg        string         `gorm:"column:req_msg;type:varchar(255);not null;comment:申请信息" json:"req_msg"`
	JoinSource    int8           `gorm:"column:join_source;type:tinyint;not null;default:2;comment:入群方式：1邀请 2申请" json:"join_source"`
	InviterUserID string         `gorm:"column:inviter_user_id;type:varchar(64);not null;comment:邀请人uid" json:"inviter_user_id"`
	HandleUserID  string         `gorm:"column:handle_user_id;type:varchar(64);not null;default:'';comment:处理人uid" json:"handle_user_id"`
	HandleTime    *time.Time     `gorm:"column:handle_time;type:timestamp;comment:处理时间" json:"handle_time,omitempty"`
	HandleResult  int8           `gorm:"column:handle_result;type:tinyint;not null;comment:1未处理 2通过 3拒绝 4取消" json:"handle_result"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamp;not null;index:idx_group_created,priority:2;index:idx_req_created,priority:2" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName GroupRequest's table name
func (GroupRequest) TableName() string {
	return "group_requests"
}
