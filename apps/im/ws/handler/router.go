package handler

import (
	"zeroIM/apps/im/ws/handler/user"
	"zeroIM/apps/im/ws/svc"
	"zeroIM/apps/im/ws/websocket"
)

func RegisterHandlers(server *websocket.Server, ctx *svc.ServiceContext) {
	server.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.Online(ctx),
		},
	})
}
