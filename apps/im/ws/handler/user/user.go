package user

import (
	"zeroIM/apps/im/ws/svc"
	"zeroIM/apps/im/ws/websocket"
)

// Online 上线
func Online(ctx *svc.ServiceContext) websocket.HandleFunc {
	return func(server *websocket.Server, conn *websocket.Conn, message *websocket.Message) {
		userIds := server.GetAllUserIds()
		uid := server.GetUid(conn)
		err := server.Send(websocket.NewMessage(uid, "", userIds), conn)
		server.Info("err ", err)
	}
}
