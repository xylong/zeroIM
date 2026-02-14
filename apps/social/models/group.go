package models

import (
	"gorm.io/gorm"
	"time"
)

// Group 群
type Group struct {
	ID              string         `gorm:"column:id;type:varchar(24);primaryKey" json:"id"`
	Name            string         `gorm:"column:name;type:varchar(255);not null;uniqueIndex:uk_creator_name;comment:群名" json:"name"`
	Icon            string         `gorm:"column:icon;type:varchar(255);not null;default:'1';comment:群图标" json:"icon"`
	Status          int8           `gorm:"column:status;type:tinyint;not null;comment:状态：1正常 2解散" json:"status"`
	CreatorUID      string         `gorm:"column:creator_uid;type:varchar(64);not null;uniqueIndex:uk_creator_name;comment:创建人uid" json:"creator_uid"`
	GroupType       int            `gorm:"column:group_type;type:int;not null;comment:群类型" json:"group_type"`
	IsVerify        int8           `gorm:"column:is_verify;type:tinyint(1);not null;comment:入群验证：1开启 2关闭" json:"is_verify"`
	Notification    string         `gorm:"column:notification;type:varchar(255);not null;default:'';comment:通知" json:"notification"`
	NotificationUID string         `gorm:"column:notification_uid;type:varchar(64);not null;default:'';comment:通知人uid" json:"notification_uid"`
	CreatedAt       time.Time      `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName Group's table name
func (Group) TableName() string {
	return "groups"
}
