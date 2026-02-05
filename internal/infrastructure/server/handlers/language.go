package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"

	"github.com/darksuei/suei-intelligence/internal/application/metadata"
	"github.com/darksuei/suei-intelligence/internal/config"
)

var Languages = []map[string]string{
	{"name": "English", "code": "EN"},
	{"name": "دری", "code": "AF"},
}

func SupportedLanguages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"languages": Languages,
	})
}

func SetLanguagePreference(c *gin.Context) {
	// Parse the request body
	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request: code is required",
		})
		return
	}

	// Validate language code against supported languages
	var isValid bool
	for _, lang := range Languages {
		if lang["code"] == req.Code {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported language code",
		})
		return
	}

	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Update metadata with the language
	if err := metadata.SetLanguage(req.Code, &databaseCfg); err != nil {
		log.Printf("Error updating language: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update language",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Language updated successfully",
		"code":    req.Code,
	})
}