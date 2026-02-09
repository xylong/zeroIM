package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DSN string
	}

	Cache []struct {
		Host string `yaml:"Host"`
		Type string `yaml:"Type"`
		Pass string `yaml:"Pass"`
	} `yaml:"Cache"`

	Jwt struct {
		AccessSecret string `yaml:"AccessSecret"`
		AccessExpire int64  `yaml:"AccessExpire"`
	} `yaml:"Jwt"`
}
