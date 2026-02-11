package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"zeroIM/apps/social/rpc/internal/config"
	"zeroIM/apps/social/rpc/internal/dao"
)

type ServiceContext struct {
	Config config.Config

	DB  *gorm.DB
	Dao *dao.Query
	Rdb *redis.Client

	SocialInfoSF syncx.SingleFlight // 全局共享
}

func NewServiceContext(c config.Config) *ServiceContext {
	var gormLogger logger.Interface = logger.Default
	if c.Mode == service.DevMode || c.Mode == service.TestMode {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	}

	db, err := gorm.Open(mysql.Open(c.Mysql.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}

	var rdb *redis.Client
	if len(c.Cache) > 0 {
		node := c.Cache[0]
		rdb = redis.NewClient(&redis.Options{
			Addr:     node.Host,
			Password: node.Pass,
			DB:       0,
		})
	}

	return &ServiceContext{
		Config:       c,
		DB:           db,
		Dao:          dao.Use(db),
		Rdb:          rdb,
		SocialInfoSF: syncx.NewSingleFlight(),
	}
}
