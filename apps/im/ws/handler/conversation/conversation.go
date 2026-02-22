package conversation

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
	"zeroIM/apps/im/ws/internal/logic"
	"zeroIM/apps/im/ws/internal/types"
	"zeroIM/apps/im/ws/svc"
	"zeroIM/apps/im/ws/websocket"
	"zeroIM/pkg/constants"
)

func Chat(svcCtx *svc.ServiceContext) websocket.HandleFunc {
	return func(server *websocket.Server, conn *websocket.Conn, message *websocket.Message) {
		var chat types.Chat
		if err := mapstructure.Decode(message.Data, &chat); err != nil {
			_ = conn.Send(websocket.NewErrMessage(err))
			return
		}

		userId := server.GetUid(conn)
		if userId == "" {
			_ = conn.Send(websocket.NewErrMessage(errors.New("user id missing")))
			return
		}

		conv := logic.NewConversation(context.Background(), server, svcCtx)

		switch chat.ChatType {
		case constants.SingleChatType:
			_ = server.SendByUserId(websocket.NewMessage(userId, chat.RecvId, chat), chat.RecvId)
			if err := conv.SingleChat(&chat, userId); err != nil {
				_ = conn.Send(websocket.NewErrMessage(err))
				return
			}
		default:
			_ = conn.Send(websocket.NewErrMessage(errors.New("unsupported chat type")))
			return
		}
	}
}
