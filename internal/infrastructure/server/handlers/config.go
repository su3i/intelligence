package handlers

import (
	"net/http"

	"github.com/darksuei/suei-intelligence/internal/config"
	"github.com/gin-gonic/gin"
)

func RetrieveConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"enforce_mfa": config.Common().EnforceMfa,
	})
	return
}