package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	accountService "github.com/darksuei/suei-intelligence/internal/application/account"
	"github.com/darksuei/suei-intelligence/internal/application/mfa"
	"github.com/darksuei/suei-intelligence/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

func RetrieveTotpURI(c *gin.Context) {
	// Parse the request body
	var req struct {
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
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

	// Retrieve account
	_account, err := accountService.RetrieveAccountWithPassword(req.Email, req.Password, &databaseCfg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _account == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid account.",
		})
		return
	}

	// Retrieve TOTP URI
	uri, err := mfa.RetrieveTotpURI(req.Email, _account.MFASecret)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve TOTP URI",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"uri": uri,
	  })
	return
}

// Confirm and Enable MFA
func ConfirmMFA(c *gin.Context) {
	// Parse the request body
	var req struct {
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Code string `json:"code" binding:"required"`
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

	// Retrieve account
	_account, err := accountService.RetrieveAccountWithPassword(req.Email, req.Password, &databaseCfg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _account == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid account.",
		})
		return
	}

	codeUint64, err := strconv.ParseUint(req.Code, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mfa code format"})
		return
	}

	code := uint32(codeUint64)

	// Confirm and enable MFA
	isCodeValid := mfa.VerifyTOTP(_account.MFASecret, code, time.Now())

	if !isCodeValid {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid TOTP code.",
		})
		return
	}

	accountService.EnableTOTP(req.Email, &databaseCfg)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	return
}