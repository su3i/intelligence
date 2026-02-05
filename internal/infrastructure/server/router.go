package server

import (
	"github.com/darksuei/suei-intelligence/internal/infrastructure/server/handlers"
	"github.com/gin-gonic/gin"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()

	// Health
	router.GET("/health", handlers.Health)

	// Language Settings
	router.GET("/supported-languages", handlers.SupportedLanguages)
	router.PUT("/set-language", handlers.SetLanguagePreference)

	// Organization
	router.POST("/organization", handlers.NewOrganization)
	router.GET("/organization/:key", handlers.RetrieveOrganization)

	// Metrics
	handlers.MetricsHandler(router)

	return router
}
