package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"

	organizationService "github.com/darksuei/suei-intelligence/internal/application/organization"
	"github.com/darksuei/suei-intelligence/internal/config"
)

func NewOrganization(c *gin.Context) {
	// Parse the request body
	var req struct {
		Name string `json:"name" binding:"required"`
		Scope string `json:"scope" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request: Missing required fields.",
		})
		return
	}

	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Create organization
	_organization, err := organizationService.NewOrganization(req.Name, "default", req.Scope, &databaseCfg)

	if err != nil {
		log.Printf("Error creating organization: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"organization": _organization,
	  })
}

func RetrieveOrganization(c *gin.Context) {
	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}
	
	key := "default" // Default organization key

	// Retrieve organization
	_organization, err := organizationService.RetrieveOrganization(key, &databaseCfg)

	if err != nil {
		log.Printf("Error retrieving organization: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _organization == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"organization": _organization,
	})
	return
}

func UpdateOrganization (c *gin.Context) {
	// Parse the request body
	var req struct {
		Name string `json:"name,omitempty"`
		Scope string `json:"scope,omitempty"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request: Missing required fields.",
		})
		return
	}

	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Update organization
	_organization, err := organizationService.UpdateOrganization(&req.Name, "default", &req.Scope, &databaseCfg)

	if err != nil {
		log.Printf("Error updating organization: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"organization": _organization,
	})
	return
}
