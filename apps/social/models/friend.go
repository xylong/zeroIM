package models

import (
	"gorm.io/gorm"
	"time"
)

// Friend 好友模型
type Friend struct {
	ID        uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    string         `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:uk_user_friend;comment:用户ID" json:"user_id"`
	FriendUID string         `gorm:"column:friend_uid;type:varchar(64);not null;uniqueIndex:uk_user_friend;index:idx_friend_uid;comment:好友UID" json:"friend_uid"`
	Remark    string         `gorm:"column:remark;type:varchar(255);not null;comment:备注" json:"remark"`
	AddSource int8           `gorm:"column:add_source;type:tinyint;not null;default:1;comment:添加方式" json:"add_source"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName Friend's table name
func (Friend) TableName() string {
	return "friends"
}
