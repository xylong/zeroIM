package models

import (
	"gorm.io/gorm"
	"time"
)

// Friend 好友关系表(单向记录，双向好友需两条记录）
type Friend struct {
	ID        int64          `gorm:"column:id;primaryKey;type:int unsigned;autoIncrement;comment:自增主键"`
	UserID    int64          `gorm:"column:user_id;type:int unsigned;not null;uniqueIndex:uk_user_friend,priority:1;index:idx_user_created,priority:1;comment:用户id"`
	FriendUID int64          `gorm:"column:friend_uid;type:int unsigned;not null;uniqueIndex:uk_user_friend,priority:2;index:idx_friend_created,priority:1;comment:好友uid"`
	Remark    string         `gorm:"column:remark;type:varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;default:'';comment:备注（对好友的备注）"`
	AddSource uint8          `gorm:"column:add_source;type:tinyint unsigned;not null;default:1;comment:添加方式：1搜索 2名片 3群聊 ..."`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime;comment:成为好友时间"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime;comment:记录更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除/拉黑时间（软删除）"`
}

// TableName Friend's table name
func (Friend) TableName() string {
	return "friends"
}
