package goodmananger

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/0LuigiCode0/hezzl/config"

	bnats "github.com/0LuigiCode0/hezzl/internal/infrastructure/broker/nats"
	rpostgres "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/postgres"
	rredis "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/redis"
	ihttp "github.com/0LuigiCode0/hezzl/internal/interfaces/http"
	"github.com/0LuigiCode0/hezzl/internal/utils"
)

func Start(ctx context.Context) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	defer utils.Shutdown()

	repoPostgres, err := rpostgres.InitRepoPostgres(ctx)
	if err != nil {
		log.Printf("ошибка инициализации postgres: %s", err)
		return
	}

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
	natsStream, err := brokerNats.ConnectStream(ctx, config.Cfg.Nats.Stream)
	if err != nil {
		log.Printf("ошибка подключения к стриму %s nats: %s", config.Cfg.Nats.Stream, err)
		return
	}

	handler := ihttp.InitHandler(repoPostgres, repoRedis, natsStream)

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := http.Serve(l, handler); err != nil {
			log.Printf("ошибка обработчика запросов: %s", err)
		}
	}()

	<-ctx.Done()
}
