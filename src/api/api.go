package api

import (
	"fmt"

	"github.com/amirazad1/creditor/api/middleware"
	"github.com/amirazad1/creditor/api/routers"
	"github.com/amirazad1/creditor/config"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	cfg := config.GetConfig()
	gin.SetMode(cfg.Server.RunMode)
	server := gin.Default()
	server.Use(gin.Logger(), gin.Recovery())
	server.Use(middleware.Cors(cfg))

	v1 := server.Group("/api/v1")
	{
		health := v1.Group("/health")
		routers.Health(health)
	}

	fmt.Println("server is running")
	server.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}
