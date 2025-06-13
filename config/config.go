package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const ServiceName = "goods-manager"

type Config struct {
	Port     int `json:"port" env:"PROXY_PORT" env-default:"90"`
	Postgres struct {
		URL string `json:"url" env:"PROXY_PG_URL" env-default:"postgres://admin:admin@localhost:5432/test"`
	} `json:"postgres"`
	ClickHouse struct {
		Addr string `json:"addr" env:"PROXY_CH_ADDR" env-default:":9440"`
		User string `json:"user" env:"PROXY_CH_USER" env-default:"admin"`
		Pwd  string `json:"pwd" env:"PROXY_CH_PWD" env-default:"admin"`
		DB   string `json:"db" env:"PROXY_CH_DB" env-default:"test"`
	}
	Nats struct {
		Addr string `json:"addr" env:"PROXY_NATS_ADDR" env-default:"nats://0.0.0.0:4222"`
		User string `json:"user" env:"PROXY_NATS_USER" env-default:"admin"`
		Pwd  string `json:"pwd" env:"PROXY_NATS_PWD" env-default:"admin"`
		Sub  string `json:"sub" env:"PROXY_NATS_SUB" env-default:"goods"`
	}
	Redis struct {
		MasterName   string `json:"master_name" env:"PROXY_REDIS_MASTER_NAME" env-default:"mymaster"`
		SentinelAddr string `json:"sentinel_addr" env:"PROXY_REDIS_SENTINEL_ADDR" env-default:":26379"`
		User         string `json:"user" env:"PROXY_REDIS_USER" env-default:"default"`
		Pwd          string `json:"pwd" env:"PROXY_REDIS_PWD" env-default:"admin"`
	}
}

// type addrs []string

// func (a *addrs) UnmarshalText(text []byte) error {
// 	addrList := strings.Split(string(text), ",")
// 	for _, v := range addrList {
// 		*a = append(*a, v)
// 	}
// 	return nil
// }

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
