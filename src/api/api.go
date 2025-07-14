package api

import (
	"fmt"

	"github.com/amirazad1/creditor/api/middleware"
	"github.com/amirazad1/creditor/api/routers"
	"github.com/amirazad1/creditor/api/validation"
	"github.com/amirazad1/creditor/config"
	"github.com/amirazad1/creditor/docs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitServer(cfg *config.Config) {
	gin.SetMode(cfg.Server.RunMode)
	server := gin.New()
	RegisterValidators()

	server.Use(middleware.DefaultStructuredLogger(cfg))
	server.Use(middleware.Cors(cfg))
	server.Use(gin.Logger(), gin.Recovery())

	RegisterRoutes(server, cfg)
	RegisterSwagger(server, cfg)

	fmt.Println("server is running")
	server.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	health := v1.Group("/health")
	routers.Health(health)

	user := v1.Group("/users")
	routers.User(user, cfg)
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

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "Creditor"
	docs.SwaggerInfo.Description = "creditor backend web api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.Port)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
