package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// ChatLog 聊天记录
type ChatLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	// 会话信息
	ConversationID string `bson:"conversation_id" json:"conversation_id"` // 会话ID（单聊=小UID_大UID，群聊=groupID）
	ChatType       int    `bson:"chat_type" json:"chat_type"`             // 1=单聊 2=群聊

	// 发送方信息
	FromUserID string `bson:"from_user_id" json:"from_user_id"`
	ToUserID   string `bson:"to_user_id,omitempty" json:"to_user_id,omitempty"` // 单聊使用
	GroupID    string `bson:"group_id,omitempty" json:"group_id,omitempty"`     // 群聊使用

	// 消息内容
	MessageID string `bson:"message_id" json:"message_id"` // 业务唯一ID（雪花ID）
	MsgType   int    `bson:"msg_type" json:"msg_type"`     // 1=文本 2=图片 3=语音 4=文件
	Content   string `bson:"content" json:"content"`       // 文本或JSON结构

	// 状态
	Status int `bson:"status" json:"status"` // 0=发送中 1=成功 2=失败

	// 扩展字段
	Ext map[string]interface{} `bson:"ext,omitempty" json:"ext,omitempty"`

	// 时间字段
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
