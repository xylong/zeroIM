package svc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"zeroIM/apps/user/rpc/internal/config"
	"zeroIM/apps/user/rpc/internal/dao"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Dao    *dao.Query
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
		Dao:    dao.Use(db),
	}
}
