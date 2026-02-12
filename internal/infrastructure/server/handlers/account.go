package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"

	accountService "github.com/darksuei/suei-intelligence/internal/application/account"
	authorizationService "github.com/darksuei/suei-intelligence/internal/application/authorization"
	"github.com/darksuei/suei-intelligence/internal/config"
	accountDomain "github.com/darksuei/suei-intelligence/internal/domain/account"
	authorizationDomain "github.com/darksuei/suei-intelligence/internal/domain/authorization"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/server/utils"
)

func NewAccount(c *gin.Context) {
	// Parse the request body
	var req struct {
		Name string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role string `json:"role" binding:"required"`
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

	role, err := accountDomain.NewAccountRole(req.Role)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid role.",
		})
		return
	}

	internalRoleKey := accountDomain.BuildRoleKey("default", authorizationDomain.AuthorizationDomainOrg, req.Role)

	internalRoleJson := map[string]string{
		"default": internalRoleKey,
	}

	// Create account
	_account, err := accountService.NewAccount(req.Name, req.Email, req.Password, role, internalRoleJson, &databaseCfg)

	if err != nil {
		log.Printf("Error creating account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"account": accountDomain.ToAccountDTO(_account),
	  })
}

func RetrieveAccounts(c *gin.Context) {
	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	allow, err := authorizationService.EnforceRoles(utils.GetUserRolesFromContext(c), "org", authorizationDomain.Organization, "read")

	if err != nil || !allow {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		return
	}

	// Retrieve accounts
	_accounts, err := accountService.RetrieveAccounts(&databaseCfg)

	if err != nil {
		log.Printf("Error retrieving account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"accounts": accountDomain.ToAccountDTOs(_accounts),
	})
}

func RetrieveAccountByEmail(c *gin.Context) {
	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Get email from query params
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required query parameter: email",
		})
		return
	}

	// Retrieve account
	_account, err := accountService.RetrieveAccount(email, &databaseCfg)

	if err != nil {
		log.Printf("Error retrieving account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _account == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"account": accountDomain.ToAccountDTO(_account),
	})
}

func UpdateAccount(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required query parameter: email",
		})
		return
	}

	var req struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
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

	_account, err := accountService.UpdateAccount(email, &req.Name, &req.Email, &databaseCfg)

	if err != nil {
		log.Printf("Error updating account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"account": _account,
	})
}