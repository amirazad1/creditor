package middleware

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/amirazad1/creditor/api/helper"
	"github.com/amirazad1/creditor/config"
	"github.com/amirazad1/creditor/pkg/limiter"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func OtpLimiter(cfg *config.Config) gin.HandlerFunc {
	var limiter = limiter.NewIPRateLimiter(rate.Every(cfg.Otp.Limiter*time.Second), 1)
	return func(c *gin.Context) {
		limiter := limiter.GetLimiter(getIP(c.Request.RemoteAddr))
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, helper.GenerateBaseResponseWithError(nil, false, helper.OtpLimiterError, errors.New("not allowed")))
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func getIP(remoteAddr string) string {
    ip, _, err := net.SplitHostPort(remoteAddr)
    if err != nil {
        return remoteAddr
    }
    return ip
}
