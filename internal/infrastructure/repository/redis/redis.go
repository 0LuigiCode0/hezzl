package rredis

import (
	"context"
	"fmt"
	"log"

	"github.com/0LuigiCode0/hezzl/config"
	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	"github.com/0LuigiCode0/hezzl/internal/utils"
	"github.com/redis/go-redis/v9"
)

type _redis struct {
	conn *redis.Client
}

type IRedis interface{}

func InitRedis(ctx context.Context) (IRedis, error) {
	redisConn := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    config.Cfg.Redis.MasterName,
		SentinelAddrs: []string{config.Cfg.Redis.SentinelAddr},
		Username:      config.Cfg.Redis.User,
		Password:      config.Cfg.Redis.Pwd,
		ClientName:    config.ServiceName,
	})

	resp := redisConn.Ping(ctx)
	if resp.Err() != nil {
		redisConn.Close()
		return nil, fmt.Errorf(consts.ErrPing, resp.Err())
	}

	utils.AddShutdown(func() {
		if err := redisConn.Close(); err != nil {
			log.Printf("ошибка закрытия соединения redis: %s", err)
		} else {
			log.Print("redis закрыт")
		}
	})

	return &_redis{conn: redisConn}, nil
}
