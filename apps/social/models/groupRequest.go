package models

import (
	"gorm.io/gorm"
	"time"
)

// GroupRequest 入群申请表
type GroupRequest struct {
	ID            int64          `gorm:"column:id;type:int(11) unsigned;primaryKey;autoIncrement:true" json:"id"`
	ReqID         string         `gorm:"column:req_id;type:varchar(64);not null;comment:请求id" json:"req_id"`                    // 请求id
	GroupID       string         `gorm:"column:group_id;type:varchar(64);not null;comment:群id" json:"group_id"`                 // 群id
	ReqMsg        string         `gorm:"column:req_msg;type:varchar(255);not null;comment:请求信息" json:"req_msg"`                 // 请求信息
	ReqTime       *time.Time     `gorm:"column:req_time;type:timestamp;not null;comment:请求时间" json:"req_time"`                  // 请求时间
	JoinSource    int64          `gorm:"column:join_source;type:tinyint(4);not null;comment:入群方式" json:"join_source"`           // 入群方式
	InviterUserID string         `gorm:"column:inviter_user_id;type:varchar(64);not null;comment:邀请人id" json:"inviter_user_id"` // 邀请人id
	HandleUserID  string         `gorm:"column:handle_user_id;type:varchar(64);not null;comment:处理人id" json:"handle_user_id"`   // 处理人id
	HandleTime    *time.Time     `gorm:"column:handle_time;type:timestamp;comment:处理时间" json:"handle_time"`                     // 处理时间
	HandleResult  int64          `gorm:"column:handle_result;type:tinyint(4);not null;comment:处理结果" json:"handle_result"`       // 处理结果
	CreatedAt     *time.Time     `gorm:"column:created_at;type:timestamp;not null;comment:创建时间" json:"created_at"`              // 创建时间
	UpdatedAt     *time.Time     `gorm:"column:updated_at;type:timestamp;not null;comment:更新时间" json:"updated_at"`              // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`                       // 删除时间
}

// TableName GroupRequest's table name
func (*GroupRequest) TableName() string {
	return "group_requests"
}
