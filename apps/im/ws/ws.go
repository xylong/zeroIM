package main

import (
	"flag"
	"fmt"
	"github.com/alicebob/miniredis/v2/server"
	"github.com/zeromicro/go-zero/core/conf"
	"zeroIM/apps/im/ws/config"
	"zeroIM/apps/im/ws/svc"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svc.NewServiceContext(c)

	srv, err := server.NewServer(c.ListenOn)
	defer srv.Close()

	// todo: 待处理

	fmt.Printf("Starting websocket server at %v ...\n", c.ListenOn)

}
