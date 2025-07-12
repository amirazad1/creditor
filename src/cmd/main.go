package main

import (
	"github.com/amirazad1/creditor/api"
	"github.com/amirazad1/creditor/config"
	"github.com/amirazad1/creditor/infra/cache"
	"github.com/amirazad1/creditor/infra/persistance/database"
	"github.com/amirazad1/creditor/infra/persistance/migration"
	"github.com/amirazad1/creditor/pkg/logging"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	err = database.InitDb(cfg)
	defer database.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	migration.Up1()

	api.InitServer(cfg)
}
