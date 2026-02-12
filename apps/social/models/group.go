package models

import (
	"gorm.io/gorm"
	"time"
)

// Group 群
type Group struct {
	ID              string         `gorm:"column:id;type:varchar(24);primaryKey" json:"id"`
	Name            string         `gorm:"column:name;type:varchar(255);not null;comment:群名" json:"name"`                               // 群名
	Icon            string         `gorm:"column:icon;type:varchar(255);not null;comment:图标" json:"icon"`                               // 图标
	Status          int64          `gorm:"column:status;type:tinyint(4);not null;comment:状态" json:"status"`                             // 状态
	CreatorUID      string         `gorm:"column:creator_uid;type:varchar(64);not null;default:0;comment:创建人用户id" json:"creator_uid"`   // 创建人用户id
	GroupType       int64          `gorm:"column:group_type;type:int(11);not null;comment:群类型" json:"group_type"`                       // 群类型
	IsVerify        int64          `gorm:"column:is_verify;type:tinyint(1);not null;comment:是否认证" json:"is_verify"`                     // 是否认证
	Notification    string         `gorm:"column:notification;type:varchar(255);not null;comment:公告通知" json:"notification"`             // 公告通知
	NotificationUID string         `gorm:"column:notification_uid;type:varchar(64);not null;comment:公告通知发布人id" json:"notification_uid"` // 公告通知发布人id
	CreatedAt       *time.Time     `gorm:"column:created_at;type:timestamp;not null;comment:创建时间" json:"created_at"`                    // 创建时间
	UpdatedAt       *time.Time     `gorm:"column:updated_at;type:timestamp;not null;comment:更新时间" json:"updated_at"`                    // 更新时间
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`                             // 删除时间
}

// TableName Group's table name
func (*Group) TableName() string {
	return "groups"
}

func (g *Group) GetVerify() bool {
	return g.IsVerify > 0
}
