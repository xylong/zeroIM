package models

import (
	"gorm.io/gorm"
	"time"
)

// Group 群
type Group struct {
	ID              int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name            string         `gorm:"column:name;type:varchar(100);not null;comment:群名" json:"name"`
	Icon            string         `gorm:"column:icon;type:varchar(255);not null;default:'';comment:群图标" json:"icon"`
	Status          int8           `gorm:"column:status;type:tinyint;not null;default:1;comment:1开启 0关闭" json:"status"`
	CreatorUID      int64          `gorm:"column:creator_uid;not null;default:0;index:idx_creator_uid;comment:创建人uid" json:"creator_uid"`
	GroupType       int8           `gorm:"column:group_type;type:tinyint;not null;default:1;comment:1=普通群 2=企业群 3=聊天室" json:"group_type"`
	IsVerify        int8           `gorm:"column:is_verify;type:tinyint;not null;comment:入群验证：0关闭 1开启" json:"is_verify"`
	Notification    string         `gorm:"column:notification;type:text;not null;comment:群公告" json:"notification"`
	NotificationUID int64          `gorm:"column:notification_uid;not null;default:0;comment:最后设置公告的人uid" json:"notification_uid"`
	MemberCount     int            `gorm:"column:member_count;type:mediumint;not null;default:1;comment:群人数" json:"member_count"`
	CreatedAt       time.Time      `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;not null;comment:更新时间" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间" json:"deleted_at"`
}

// TableName Group's table name
func (Group) TableName() string {
	return "groups"
}
