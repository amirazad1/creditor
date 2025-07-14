package routers

import (
	"github.com/amirazad1/creditor/api/handler"
	"github.com/gin-gonic/gin"
)

func Health(r *gin.RouterGroup){
	health:=handler.NewHealthHandler()
	r.GET("/",health.Check)
}