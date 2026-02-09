package models

import "time"

// User 用户
type User struct {
	ID        string     `gorm:"column:id;type:varchar(24);primaryKey" json:"id"`
	Avatar    string     `gorm:"column:avatar;type:varchar(191);not null;default:''" json:"avatar"`
	Nickname  string     `gorm:"column:nickname;type:varchar(24);not null" json:"nickname"`
	Phone     string     `gorm:"column:phone;type:varchar(20);not null" json:"phone"`
	Password  *string    `gorm:"column:password;type:varchar(191)" json:"password,omitempty"`
	Status    *int8      `gorm:"column:status;type:tinyint" json:"status"`
	Sex       *int8      `gorm:"column:sex;type:tinyint" json:"sex"`
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetPassword() string {
	if u.Password == nil {
		return ""
	}

	return *u.Password
}
