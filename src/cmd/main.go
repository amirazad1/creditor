package main

import (
	"log"

	"github.com/amirazad1/creditor/api"
	"github.com/amirazad1/creditor/config"
	"github.com/amirazad1/creditor/infra/cache"
	"github.com/amirazad1/creditor/infra/persistance/database"
)

func main() {
	cfg := config.GetConfig()

	if err := cache.InitRedis(cfg); err != nil {
		log.Fatalf("Redis initialize error: %s", err)
	}
	defer cache.CloseRedis()

	if err := database.InitDb(cfg); err != nil {
		log.Fatalf("Postgres initialize error: %s", err)
	}
	defer database.CloseDb()

	api.InitServer(cfg)
}
