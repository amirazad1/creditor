package routers

import (
	"github.com/amirazad1/creditor/api/handlers"
	"github.com/gin-gonic/gin"
)

func Health(r *gin.RouterGroup){
	health:=handlers.NewHealthHandler()
	r.GET("/",health.Check)
}