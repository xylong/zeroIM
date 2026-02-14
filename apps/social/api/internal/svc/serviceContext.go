package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zeroIM/apps/social/api/internal/config"
	"zeroIM/apps/social/rpc/socialClient"
	"zeroIM/apps/user/rpc/userClient"
)

type ServiceContext struct {
	Config config.Config

	// rpc
	User   userClient.User
	Social socialClient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		User:   userClient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social: socialClient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
}
