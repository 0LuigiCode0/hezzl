package config

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServiceName string `json:"service_name" env:"PROXY_SERVICE_NAME" env-default:"goods-manager"`
	Port        int    `json:"port" env:"PROXY_PORT" env-default:"90"`

	Postgres struct {
		URL string `json:"url" env:"PROXY_PG_URL" env-default:"postgres://admin:admin@host.docker.internal:5432/test"`
	} `json:"postgres"`

	ClickHouse struct {
		Addr string `json:"addr" env:"PROXY_CH_ADDR" env-default:"host.docker.internal:9000"`
		User string `json:"user" env:"PROXY_CH_USER" env-default:"admin"`
		Pwd  string `json:"pwd" env:"PROXY_CH_PWD" env-default:"admin"`
		DB   string `json:"db" env:"PROXY_CH_DB" env-default:"test"`
	}

	Nats struct {
		Addr string `json:"addr" env:"PROXY_NATS_ADDR" env-default:"nats://host.docker.internal:4222"`
		User string `json:"user" env:"PROXY_NATS_USER" env-default:"admin"`
		Pwd  string `json:"pwd" env:"PROXY_NATS_PWD" env-default:"admin"`

		Subj   string `json:"subj" env:"PROXY_NATS_SUB" env-default:"log.goods"`
		Stream string `json:"stream" env:"PROXY_NATS_STREAM" env-default:"logs"`

		Batch      int           `json:"batch" env:"PROXY_NATS_BATCH" env-default:"10"`
		Delay      time.Duration `json:"delay" env:"PROXY_NATS_DELAY" env-default:"10s"`
		RetryCount int           `json:"retry_count" env:"PROXY_NATS_RETRY_COUNT" env-default:"4"`
	}

	Redis struct {
		MasterName   string `json:"master_name" env:"PROXY_REDIS_MASTER_NAME" env-default:"mymaster"`
		SentinelAddr string `json:"sentinel_addr" env:"PROXY_REDIS_SENTINEL_ADDR" env-default:"host.docker.internal:26379"`
		User         string `json:"user" env:"PROXY_REDIS_USER" env-default:"default"`
		Pwd          string `json:"pwd" env:"PROXY_REDIS_PWD" env-default:"admin"`

		RetryCount int           `json:"retry_count" env:"PROXY_REDIS_RETRY_COUNT" env-default:"4"`
		Expire     time.Duration `json:"expire" env:"PROXY_REDIS_EXPIRE" env-default:"1m"`
	}
}

var Cfg *Config

func ParseConfig(file string) error {
	Cfg = new(Config)
	f, err := os.Open(file)
	if err == nil {
		buf, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		return json.Unmarshal(buf, Cfg)
	}

	return cleanenv.ReadEnv(Cfg)
}
