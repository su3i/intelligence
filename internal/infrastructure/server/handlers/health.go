package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	log.Print("Received health check request.")

	c.JSON(http.StatusOK, gin.H{
		"message": "Healthy",
		"version": "v0.0.0",
		"copyright": "2026 SUEI INC.",
		"uptime": "",
	  })
}