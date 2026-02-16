package models

import (
	"gorm.io/gorm"
	"time"
)

// Group 群
type Group struct {
	ID              string         `gorm:"column:id;type:varchar(24);primaryKey"`
	Name            string         `gorm:"column:name;type:varchar(255);not null;uniqueIndex:uk_creator_name"`
	Icon            string         `gorm:"column:icon;type:varchar(255);not null;default:1;comment:群图标"`
	Status          int8           `gorm:"column:status;type:tinyint(4);not null;default:1;comment:状态：1正常 2解散"`
	CreatorUID      string         `gorm:"column:creator_uid;type:varchar(64);not null;uniqueIndex:uk_creator_name;comment:创建人uid"`
	GroupType       int8           `gorm:"column:group_type;type:tinyint(1);not null;default:1;comment:群类型"`
	IsVerify        int8           `gorm:"column:is_verify;type:tinyint(1);not null;comment:入群验证：1开启 2关闭"`
	Notification    string         `gorm:"column:notification;type:varchar(255);not null;default:'';comment:通知"`
	NotificationUID string         `gorm:"column:notification_uid;type:varchar(64);not null;default:'';comment:通知人uid"`
	CreatedAt       time.Time      `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName Group's table name
func (Group) TableName() string {
	return "groups"
}
