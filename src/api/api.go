package api

import (
	"fmt"
	"github.com/amirazad1/creditor/api/routers"
	"github.com/gin-gonic/gin"
)

func InitServer(){
	server:=gin.Default()
	server.Use(gin.Logger(),gin.Recovery())

	v1:=server.Group("/api/v1")
	{
		health:=v1.Group("/health")
		routers.Health(health)
	}
	fmt.Println("server is running")
	server.Run(":8090")	
}