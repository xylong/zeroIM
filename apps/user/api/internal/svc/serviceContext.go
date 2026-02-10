package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zeroIM/apps/user/api/internal/config"
	"zeroIM/apps/user/rpc/userClient"
)

type ServiceContext struct {
	Config config.Config

	userClient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		User: userClient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
