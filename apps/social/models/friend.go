package models

import (
	"gorm.io/gorm"
	"time"
)

// Friend 好友模型
type Friend struct {
	Id        int64          `gorm:"column:id;type:int(11) unsigned;primaryKey;autoIncrement:true" json:"id"`
	UserId    string         `gorm:"column:user_id;type:varchar(64);not null;comment:用户id" json:"user_id"`         // 用户id
	FriendUid string         `gorm:"column:friend_uid;type:varchar(64);not null;comment:好友用户id" json:"friend_uid"` // 好友用户id
	Remark    string         `gorm:"column:remark;type:varchar(255);not null;comment:备注" json:"remark"`            // 备注
	AddSource int64          `gorm:"column:add_source;type:tinyint(4);not null;comment:添加渠道" json:"add_source"`    // 添加渠道
	CreatedAt *time.Time     `gorm:"column:created_at;type:timestamp;not null;comment:创建时间" json:"created_at"`     // 创建时间
	UpdatedAt *time.Time     `gorm:"column:updated_at;type:timestamp;not null;comment:删除时间" json:"updated_at"`     // 删除时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`              // 删除时间
}

// TableName Friend's table name
func (*Friend) TableName() string {
	return "friends"
}
