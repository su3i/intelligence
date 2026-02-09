package server

import (
	"github.com/darksuei/suei-intelligence/internal/infrastructure/server/handlers"
	middleware "github.com/darksuei/suei-intelligence/internal/infrastructure/server/middlewares"
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
	router.GET("/organization", middleware.AuthMiddleware(), handlers.RetrieveOrganization)

	// Account
	router.POST("/account", handlers.NewAccount)
	router.GET("/account", handlers.RetrieveAccountByEmail)
	router.GET("/accounts", middleware.AuthMiddleware(), handlers.RetrieveAccounts)

	// MFA
	router.POST("/mfa/totp-uri", handlers.RetrieveTotpURI)
	router.POST("/mfa/confirm", handlers.ConfirmMFA)

	// Auth
	router.POST("/auth/login", handlers.Login)
	router.POST("/auth/mfa", handlers.MFA)
	router.POST("/auth/refresh-token", handlers.RefreshToken)
	router.POST("/auth/revoke-token", handlers.RevokeToken)

	// Project
	router.POST("/project", middleware.AuthMiddleware(), handlers.NewProject)
	router.GET("/project/:key", middleware.AuthMiddleware(), handlers.RetrieveProject)
	router.PUT("/project/:key", middleware.AuthMiddleware(), handlers.UpdateProject)
	router.GET("/projects", middleware.AuthMiddleware(), handlers.RetrieveProjects)

	// Metrics
	handlers.MetricsHandler(router)

	return router
}
