package models

import (
	"gorm.io/gorm"
	"time"
)

// GroupMember 群成员表
type GroupMember struct {
	ID          uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GroupID     string         `gorm:"column:group_id;type:varchar(64);not null;uniqueIndex:uk_group_user;comment:群id" json:"group_id"`
	UserID      string         `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:uk_group_user;comment:用户uid" json:"user_id"`
	RoleLevel   int8           `gorm:"column:role_level;type:tinyint;not null;default:3;comment:1.创建者 2.管理者 3.普通" json:"role_level"`
	JoinTime    *time.Time     `gorm:"column:join_time;type:timestamp;comment:入群时间" json:"join_time,omitempty"`
	JoinSource  int8           `gorm:"column:join_source;type:tinyint;default:1;comment:入群方式：1.邀请，2.申请" json:"join_source"`
	InviterUID  string         `gorm:"column:inviter_uid;type:varchar(64);default:'';comment:邀请人uid" json:"inviter_uid"`
	OperatorUID string         `gorm:"column:operator_uid;type:varchar(64);default:'';comment:操作人uid" json:"operator_uid"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName GroupMember's table name
func (GroupMember) TableName() string {
	return "group_members"
}
