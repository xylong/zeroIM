package models

import (
	"gorm.io/gorm"
	"time"
)

// FriendRequest 好友申请记录表（单向申请）
type FriendRequest struct {
	ID           int64          `gorm:"column:id;primaryKey;type:int unsigned;autoIncrement;comment:自增主键"`
	UserID       int64          `gorm:"column:user_id;type:int unsigned;not null;uniqueIndex:uk_user_target,priority:1;index:idx_user_created,priority:1;comment:发起申请的用户id"`
	ReqUID       int64          `gorm:"column:req_uid;type:int unsigned;not null;uniqueIndex:uk_user_target,priority:2;index:idx_target_status_time,priority:1;comment:被申请的好友uid（目标用户）"`
	ReqMsg       string         `gorm:"column:req_msg;type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;default:'';comment:申请附言/验证消息"`
	HandleResult uint8          `gorm:"column:handle_result;type:tinyint unsigned;not null;default:0;index:idx_target_status_time,priority:2;comment:0=待处理 1=通过 2=拒绝 3=申请人取消"`
	HandleMsg    string         `gorm:"column:handle_msg;type:varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;default:'';comment:处理时的回复/拒绝理由"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime;comment:申请创建时间"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime;comment:记录更新时间"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index;comment:软删除时间"`
}

// TableName FriendRequest's table name
func (FriendRequest) TableName() string {
	return "friend_requests"
}
