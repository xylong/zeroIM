package models

import (
	"gorm.io/gorm"
	"time"
)

// GroupMember 群成员表
type GroupMember struct {
	ID          int64          `gorm:"column:id;type:int(11) unsigned;primaryKey;autoIncrement:true" json:"id"`
	GroupID     string         `gorm:"column:group_id;type:varchar(64);not null;comment:群id" json:"group_id"`                   // 群id
	UserID      string         `gorm:"column:user_id;type:varchar(64);not null;comment:用户id" json:"user_id"`                    // 用户id
	RoleLevel   int64          `gorm:"column:role_level;type:tinyint(4);not null;comment:角色等级" json:"role_level"`               // 角色等级
	JoinTime    time.Time      `gorm:"column:join_time;type:timestamp;not null;comment:加入时间" json:"join_time"`                  // 加入时间
	JoinSource  int64          `gorm:"column:join_source;type:tinyint(4);not null;comment:加入方式" json:"join_source"`             // 加入方式
	InviterUID  string         `gorm:"column:inviter_uid;type:varchar(64);not null;default:0;comment:邀请人id" json:"inviter_uid"` // 邀请人id
	OperatorUID string         `gorm:"column:operator_uid;type:varchar(64);not null;comment:操作人id" json:"operator_uid"`         // 操作人id
	CreatedAt   *time.Time     `gorm:"column:created_at;type:timestamp;not null;comment:创建时间" json:"created_at"`                // 创建时间
	UpdatedAt   *time.Time     `gorm:"column:updated_at;type:timestamp;not null;comment:更新时间" json:"updated_at"`                // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`                         // 删除时间
}

// TableName GroupMember's table name
func (*GroupMember) TableName() string {
	return "group_members"
}
