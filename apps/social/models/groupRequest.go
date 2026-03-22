package models

import (
	"gorm.io/gorm"
	"time"
)

// GroupRequest 入群申请表
type GroupRequest struct {
	ID            int64          `gorm:"primaryKey;column:id;type:int unsigned;autoIncrement;comment:自增主键"`
	ReqID         int64          `gorm:"column:req_id;type:int unsigned;not null;index:idx_requester_status;comment:申请人/被邀请人 uid"`
	GroupID       int64          `gorm:"column:group_id;type:int unsigned;not null;index:idx_group_pending,priority:1;index:idx_group_req,priority:1;comment:群id"`
	ReqMsg        string         `gorm:"column:req_msg;type:varchar(100);not null;default:'';comment:申请/邀请附言"`
	JoinSource    uint8          `gorm:"column:join_source;type:tinyint unsigned;not null;default:2;comment:1=被邀请入群 2=主动申请"`
	InviterUserID int64          `gorm:"column:inviter_user_id;type:int unsigned;not null;default:0;comment:邀请人uid（join_source=1时有效）"`
	HandleUserID  int64          `gorm:"column:handle_user_id;type:int unsigned;not null;default:0;comment:处理人uid（群管理员/群主）"`
	HandleResult  uint8          `gorm:"column:handle_result;type:tinyint unsigned;not null;default:0;index:idx_group_pending,priority:2;index:idx_requester_status,priority:2;comment:0=待处理 1=通过 2=拒绝 3=取消"`
	HandledAt     *time.Time     `gorm:"column:handled_at;comment:处理时间"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间（软删除）"`
}

// TableName GroupRequest's table name
func (GroupRequest) TableName() string {
	return "group_requests"
}
