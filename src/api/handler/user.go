package handler

import (
	"net/http"

	"github.com/amirazad1/creditor/api/dto"
	"github.com/amirazad1/creditor/api/helper"
	"github.com/amirazad1/creditor/config"
	"github.com/amirazad1/creditor/service"
	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	otpService *service.OtpService
	config     *config.Config
}

func NewUserHandler(cfg *config.Config) *UsersHandler {
	otpService := service.NewOtpService(cfg)
	return &UsersHandler{otpService: otpService, config: cfg}
}

// SendOtp godoc
// @Summary Send otp to user
// @Description Send otp to user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param Request body dto.GetOtpRequest true "GetOtpRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/send-otp [post]
func (h *UsersHandler) SendOtp(c *gin.Context) {
	req := new(dto.GetOtpRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	err = h.otpService.SendOtp(req.MobileNumber)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	// TODO: Call internal SMS service
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, helper.Success))
}
