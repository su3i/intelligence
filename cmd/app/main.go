package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darksuei/suei-intelligence/internal/config"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/database"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/server"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load env: %v", err)
	}

	var commonCfg config.CommonConfig
	err = envconfig.Process("", &commonCfg)
	if err != nil {
		log.Fatalf("Failed to load common config: %v", err)
	}

	var databaseCfg config.DatabaseConfig
	err = envconfig.Process("", &databaseCfg)
	if err != nil {
		log.Fatalf("Failed to load common config: %v", err)
	}

	database.Initialize(&databaseCfg)
	database.Migrate(&databaseCfg)

	router := server.InitializeRouter()

	httpServer := &http.Server{
		Addr:    ":" + commonCfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Printf("Application is running on port: %s", commonCfg.AppPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Application startup failed: %s", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down application..")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %s", err)
	}

	log.Println("Application shutdown successfully..")
}
