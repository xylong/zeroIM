package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DSN string
	}

	Cache []struct {
		Host string
		Type string
		Pass string
	}
}
