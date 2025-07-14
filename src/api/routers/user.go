package routers

import (
	"github.com/amirazad1/creditor/api/handler"
	"github.com/amirazad1/creditor/api/middleware"
	"github.com/amirazad1/creditor/config"
	"github.com/gin-gonic/gin"
)

func User(u *gin.RouterGroup, cfg *config.Config) {
	user := handler.NewUserHandler(cfg)
	u.POST("/send-otp", middleware.OtpLimiter(cfg), user.SendOtp)
}
