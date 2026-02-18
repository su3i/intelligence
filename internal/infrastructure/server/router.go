package server

import (
	"github.com/darksuei/suei-intelligence/internal/infrastructure/server/handlers"
	middleware "github.com/darksuei/suei-intelligence/internal/infrastructure/server/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()

	// Cors Settings
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
	}))

	// Health
	router.GET("/health", handlers.Health)

	// Config
	router.GET("/config", handlers.RetrieveConfig)

	// Language Settings
	router.GET("/supported-languages", handlers.SupportedLanguages)
	router.PUT("/set-language", handlers.SetLanguagePreference)
	router.GET("/get-language", handlers.RetrieveLanguagePreference)

	// Organization
	router.POST("/organization", handlers.NewOrganization)
	router.PUT("/organization", handlers.UpdateOrganization)
	router.GET("/organization", handlers.RetrieveOrganization)

	// Account
	router.POST("/account", handlers.NewAccount)
	router.GET("/account", handlers.RetrieveAccountByEmail)
	router.PUT("/account", handlers.UpdateAccount)
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

	// Datasource
	router.GET("/supported-datasources", middleware.AuthMiddleware(), handlers.SupportedDatasources)
	router.GET("/supported-datasources/:sourceType", middleware.AuthMiddleware(), handlers.SupportedDatasource)
	router.POST("/project/:key/datasources", middleware.AuthMiddleware(), handlers.NewDatasource)
	router.GET("/project/:key/datasources", middleware.AuthMiddleware(), handlers.RetrieveDatasources)
	router.DELETE("/project/:key/datasources/:id", middleware.AuthMiddleware(), handlers.DeleteDatasource)

	// Datasource - schemas
	router.GET("/project/:key/datasources/:id/source-schema-def", middleware.AuthMiddleware(), handlers.RetrieveSourceSchema)

	// Metrics
	handlers.MetricsHandler(router)

	return router
}
