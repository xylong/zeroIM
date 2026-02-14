package models

import (
	"gorm.io/gorm"
	"time"
)

// User 用户
type User struct {
	ID        string         `gorm:"column:id;type:varchar(24);primaryKey" json:"id"`
	Avatar    string         `gorm:"column:avatar;type:varchar(191);not null;default:'';comment:头像" json:"avatar"`
	Nickname  string         `gorm:"column:nickname;type:varchar(24);not null;default:'';comment:昵称" json:"nickname"`
	Phone     string         `gorm:"column:phone;type:varchar(20);not null;uniqueIndex:uk_phone;comment:手机号" json:"phone"`
	Password  string         `gorm:"column:password;type:varchar(191);not null;default:'';comment:密码" json:"-"`
	Status    int8           `gorm:"column:status;type:tinyint;not null;default:1;comment:状态 1正常 2禁用" json:"status"`
	Sex       int8           `gorm:"column:sex;type:tinyint;not null;default:0;comment:性别 1男 2女 3未知" json:"sex"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (u User) TableName() string {
	return "users"
}
