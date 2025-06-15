package logmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/0LuigiCode0/hezzl/config"
	dclickhouse "github.com/0LuigiCode0/hezzl/internal/domain/clickhouse"
	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
	bnats "github.com/0LuigiCode0/hezzl/internal/infrastructure/broker/nats"
	rclickhouse "github.com/0LuigiCode0/hezzl/internal/infrastructure/repository/clickhouse"
	"github.com/0LuigiCode0/hezzl/internal/utils"
	"github.com/nats-io/nats.go"
)

func Start(ctx context.Context) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	defer utils.Shutdown()

	repoClickHouse, err := rclickhouse.InitClickHouse(ctx)
	if err != nil {
		log.Printf("ошибка инициализации clickhouse: %s", err)
		return
	}

	brokerNats, err := bnats.InitNats(ctx)
	if err != nil {
		log.Printf("ошибка инициализации nats: %s", err)
		return
	}
	natsStream, err := brokerNats.CreateStream(ctx, config.Cfg.Nats.Stream, "log.goods")
	if err != nil {
		log.Printf("ошибка создания/подключения к стриму %s nats: %s", config.Cfg.Nats.Stream, err)
		return
	}
	consumerNats, err := natsStream.Subscribe(ctx, "log.goods", config.Cfg.ServiceName)
	if err != nil {
		log.Printf("ошибка создания подписчика nats: %s", err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := logsStore(ctx, repoClickHouse, consumerNats); err != nil {
			log.Printf("ошибка обработчика логов: %s", err)
		}
	}()

	<-ctx.Done()
}

func logsStore(ctx context.Context, repoClickHouse rclickhouse.IClickHouse, consumerNats bnats.IConsumer) error {
	msgProgress := make([]*nats.Msg, 0, config.Cfg.Nats.Batch)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(config.Cfg.Nats.Delay):
			msgs, err := consumerNats.Fetch(ctx, config.Cfg.Nats.Batch)
			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					log.Print(err)
				}
				continue
			}

			batch, err := repoClickHouse.InsertGoodsLogBatch(ctx)
			if err != nil {
				log.Print(err)
				continue
			}

			msgProgress = msgProgress[:0]
			for _, msg := range msgs {
				eventLog := new(dclickhouse.LogEventGood)
				err = json.Unmarshal(msg.Data, eventLog)
				if err != nil {
					log.Print(fmt.Errorf(consts.ErrJsonUnmarshal, err))
					msg.Nak()
					continue
				}

				err = batch.AppendStruct(eventLog)
				if err != nil {
					log.Print(err)
					msg.Nak()
					continue
				}

				msgProgress = append(msgProgress, msg)
			}

			var isNak bool
			err = batch.Send()
			if err != nil {
				isNak = true
				log.Print(err)
				continue
			}
			for _, v := range msgProgress {
				if isNak {
					v.Nak()
				} else {
					v.Ack()
				}
			}
		}
	}
}
