package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"zeroIM/apps/im/ws/handler"
	"zeroIM/apps/im/ws/internal/config"
	"zeroIM/apps/im/ws/internal/svc"
	"zeroIM/apps/im/ws/internal/websocket"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	//ctx := svc.NewServiceContext(c)
	//srv := websocket.NewServer(c.ListenOn,
	//	websocket.WithAuthentication(handler.NewJwtAuth(ctx)),
	//)
	//defer srv.Stop()
	//
	//handler.RegisterHandlers(srv, ctx)
	//
	//fmt.Printf("Starting websocket server at %v ...\n", c.ListenOn)
	//srv.Start()

	ctx := svc.NewServiceContext(c)
	s := websocket.NewServer(websocket.WithAuthentication(handler.NewJwtAuth(ctx)))
	s.Run()
}
