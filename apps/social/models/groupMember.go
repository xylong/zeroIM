package models

import (
	"gorm.io/gorm"
	"time"
)

// GroupMember 群成员关系表
type GroupMember struct {
	ID              int64          `gorm:"column:id;primaryKey;type:int unsigned;autoIncrement;comment:自增主键"`
	GroupID         int64          `gorm:"column:group_id;type:int unsigned;not null;uniqueIndex:uk_group_user,priority:1;index:idx_group_role,priority:1;comment:群id"`
	UserID          int64          `gorm:"column:user_id;type:int unsigned;not null;uniqueIndex:uk_group_user,priority:2;index:idx_user_group,priority:1;comment:用户uid"`
	RoleLevel       uint8          `gorm:"column:role_level;type:tinyint unsigned;not null;default:3;index:idx_group_role,priority:2;comment:0普通成员，10管理员，20群主"`
	JoinTime        *time.Time     `gorm:"column:join_time;comment:入群时间"`
	JoinSource      uint8          `gorm:"column:join_source;type:tinyint unsigned;not null;default:2;comment:1=被邀请 2=主动申请通过"`
	InviterUID      int64          `gorm:"column:inviter_uid;type:int unsigned;not null;default:0;comment:邀请人uid（join_source=1时有效）"`
	LastOperatorUID int64          `gorm:"column:last_operator_uid;type:int unsigned;not null;default:0;comment:最后操作人uid（添加/修改/移除）"`
	Status          int8           `gorm:"column:status;type:tinyint;not null;default:1;comment:1正常，2禁言，3被踢"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime;comment:记录创建时间"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime;comment:记录更新时间"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index;uniqueIndex:uk_group_user,priority:3;comment:退出/被踢时间（软删除）"`
}

// TableName GroupMember's table name
func (GroupMember) TableName() string {
	return "group_members"
}
