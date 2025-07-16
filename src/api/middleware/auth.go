package middleware

import (
	"net/http"
	"strings"

	"github.com/amirazad1/creditor/api/helper"
	"github.com/amirazad1/creditor/config"
	constant "github.com/amirazad1/creditor/constant"
	"github.com/amirazad1/creditor/pkg/service_errors"
	"github.com/amirazad1/creditor/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	var tokenService = service.NewTokenService(cfg)

	return func(c *gin.Context) {
		var err error
		claimMap := map[string]interface{}{}
		auth := c.GetHeader(constant.AuthorizationHeaderKey)
		token := strings.Split(auth, " ")
		if auth == "" || len(token) < 2 {
			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired}
		} else {
			claimMap, err = tokenService.GetClaims(token[1])
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
				default:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
				}
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.GenerateBaseResponseWithError(
				nil, false, helper.AuthError, err,
			))
			return
		}

		c.Set(constant.UserIdKey, claimMap[constant.UserIdKey])
		c.Set(constant.FirstNameKey, claimMap[constant.FirstNameKey])
		c.Set(constant.LastNameKey, claimMap[constant.LastNameKey])
		c.Set(constant.RolesKey, claimMap[constant.RolesKey])
		c.Set(constant.ExpireTimeKey, claimMap[constant.ExpireTimeKey])

		c.Next()
	}
}
