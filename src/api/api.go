package api

import (
	"fmt"

	"github.com/amirazad1/creditor/api/middleware"
	"github.com/amirazad1/creditor/api/routers"
	"github.com/amirazad1/creditor/api/validation"
	"github.com/amirazad1/creditor/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitServer() {
	cfg := config.GetConfig()
	gin.SetMode(cfg.Server.RunMode)
	server := gin.New()
	RegisterValidators()
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

func RegisterValidators() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", validation.IranianMobileNumberValidator, true)
		if err != nil {
			fmt.Errorf("%s", err.Error())
		}
		err = val.RegisterValidation("password", validation.PasswordValidator, true)
		if err != nil {
			fmt.Errorf("%s", err.Error())
		}
	}
}
