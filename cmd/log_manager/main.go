package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/0LuigiCode0/hezzl/config"
	logmanager "github.com/0LuigiCode0/hezzl/internal/infrastructure/service/log_manager"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "cfg", "./config-server.json", "")
	flag.Parse()

	err := config.ParseConfig(filePath)
	if err != nil {
		log.Printf("ошибка чтения конфига: %s", err)
		return
	}

	ctx, close := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer close()

	logmanager.Start(ctx)
	log.Print("сервер остановлен")
}
