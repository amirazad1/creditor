package main

import (
	"fmt"

	"github.com/amirazad1/creditor/api"
	"github.com/amirazad1/creditor/config"
	"github.com/amirazad1/creditor/infra/cache"
)

func main() {
	cfg := config.GetConfig()
	err := cache.InitRedis(cfg)
	fmt.Printf("Error in initialize redis: %s", err.Error())
	defer cache.CloseRedis()
	api.InitServer(cfg)
}
