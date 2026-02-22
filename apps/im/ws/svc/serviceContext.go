package svc

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"zeroIM/apps/im/ws/internal/config"
)

const (
	mongoConnectTimeout = 10 * time.Second
	chatLogColl         = "chat_log"
)

type ServiceContext struct {
	Config config.Config

	MongoClient *mongo.Client
	MongoDB     *mongo.Database
	ChatLogColl *mongo.Collection
}

func NewServiceContext(c config.Config) *ServiceContext {
	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.Mongo.Url))
	if err != nil {
		panic("mongo connect failed: " + err.Error())
	}
	if err = client.Ping(ctx, nil); err != nil {
		panic("mongo ping failed: " + err.Error())
	}

	db := client.Database(c.Mongo.Db)

	return &ServiceContext{
		Config:      c,
		MongoClient: client,
		MongoDB:     db,
		ChatLogColl: db.Collection(chatLogColl),
	}
}
