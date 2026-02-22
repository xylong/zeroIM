package logic

import (
	"context"
	"time"

	"zeroIM/apps/im/models"
	"zeroIM/apps/im/ws/internal/types"
	"zeroIM/apps/im/ws/svc"
	"zeroIM/apps/im/ws/websocket"
	"zeroIM/pkg/snowflake"
	"zeroIM/pkg/wuid"
)

// Conversation 会话
type Conversation struct {
	ctx    context.Context
	server *websocket.Server
	svcCtx *svc.ServiceContext
}

func NewConversation(ctx context.Context, server *websocket.Server, svcCtx *svc.ServiceContext) *Conversation {
	return &Conversation{ctx: ctx, server: server, svcCtx: svcCtx}
}

// SingleChat 单聊
func (c *Conversation) SingleChat(chat *types.Chat, userId string) error {
	if chat.ConversationId == "" {
		chat.ConversationId = wuid.CombineId(userId, chat.RecvId)
	}

	now := time.Now()
	chatLog := &models.ChatLog{
		ConversationID: chat.ConversationId,
		ChatType:       int(chat.ChatType),
		FromUserID:     userId,
		ToUserID:       chat.RecvId,
		MessageID:      snowflake.GenIDStr(),
		MsgType:        int(chat.MType),
		Content:        chat.Msg.Content,
		Status:         1, // 成功
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	_, err := c.svcCtx.ChatLogColl.InsertOne(c.ctx, chatLog)
	if err != nil {
		return err
	}

	// 若对方在线则推送
	//msg := websocket.NewMessage(userId, chat.RecvId, chat)
	//msg.Method = "conversation.push"
	//_ = c.server.SendByUserId(msg, chat.RecvId)

	return nil
}
