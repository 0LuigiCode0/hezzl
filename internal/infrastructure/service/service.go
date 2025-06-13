package service

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/0LuigiCode0/hezzl/config"

	bnats "github.com/0LuigiCode0/hezzl/internal/infrastructure/broker/nats"
	rpostgres "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/postgres"
	rredis "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/redis"
	ihttp "github.com/0LuigiCode0/hezzl/internal/interfaces/http"
	"github.com/0LuigiCode0/hezzl/internal/utils"
)

func Start(ctx context.Context) {
	defer utils.Shutdown()

	repoPostgres, err := rpostgres.InitRepoPostgres(ctx)
	if err != nil {
		log.Printf("ошибка инициализации postgres: %s", err)
		return
	}

	// repoClickHouse, err := repoclickhouse.InitClickHouse(ctx)
	// if err != nil {
	// 	log.Printf("ошибка инициализации clickhouse: %w", err)
	// 	return
	// }

	repoRedis, err := rredis.InitRedis(ctx)
	if err != nil {
		log.Printf("ошибка инициализации redis: %s", err)
		return
	}

	brokerNats, err := bnats.InitNats(ctx)
	if err != nil {
		log.Printf("ошибка инициализации nats: %s", err)
		return
	}

	handler := ihttp.InitHandler(repoPostgres, repoRedis, brokerNats)

	l, err := net.ListenTCP("tcp4", &net.TCPAddr{Port: config.Cfg.Port})
	if err != nil {
		log.Printf("ошибка инициализации слушателя на порту %d: %s", config.Cfg.Port, err)
		return
	}
	utils.AddShutdown(func() {
		if err := l.Close(); err != nil {
			log.Printf("ошибка закрытия слушателя: %s", err)
		} else {
			log.Printf("слушатель порта %d закрыт", config.Cfg.Port)
		}
	})

	go func() {
		if err := http.Serve(l, handler); err != nil {
			log.Printf("ошибка обработчика запросов: %s", err)
		}
	}()

	<-ctx.Done()
}
