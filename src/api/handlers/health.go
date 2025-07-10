package handlers

import (
	"net/http"

	"github.com/amirazad1/creditor/api/helper"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, helper.GenerateBaseResponse("The server is OK", true, helper.Success))
}
