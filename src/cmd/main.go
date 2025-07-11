package main

import (	
	"github.com/amirazad1/creditor/api"
	"github.com/amirazad1/creditor/config"
	"github.com/amirazad1/creditor/infra/cache"
	"github.com/amirazad1/creditor/infra/persistance/database"
	"github.com/amirazad1/creditor/pkg/logging"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	if err := cache.InitRedis(cfg); err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}
	defer cache.CloseRedis()

	if err := database.InitDb(cfg); err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	defer database.CloseDb()

	api.InitServer(cfg)
}
