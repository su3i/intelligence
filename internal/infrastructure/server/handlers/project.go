package handlers

import (
	"log"
	"net/http"

	authorizationService "github.com/darksuei/suei-intelligence/internal/application/authorization"
	"github.com/darksuei/suei-intelligence/internal/application/project"
	"github.com/darksuei/suei-intelligence/internal/config"
	authorizationDomain "github.com/darksuei/suei-intelligence/internal/domain/authorization"
	projectDomain "github.com/darksuei/suei-intelligence/internal/domain/project"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/server/utils"
	"github.com/gin-gonic/gin"
)

func NewProject(c *gin.Context) {
	// Parse the request body
	var req struct {
		Name string `json:"name" binding:"required"`
		Key string `json:"key" binding:"required"`
		Stage string `json:"stage" binding:"required"`
		BusinessDomain string `json:"businessDomain" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Validation failed.",
			"errors": utils.FormatValidationErrors(err),
		})
		return
	}

	createdByEmail, err := utils.GetUserEmailFromContext(c)

	if err != nil || createdByEmail == nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get account",
		})
		return
	}

	// Authorization
	allow, err := authorizationService.EnforceRoles(utils.GetUserRolesFromContext(c), "org", authorizationDomain.Organization, "write")

	if err != nil || !allow {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		return
	}

	// Create project
	_project, err := project.NewProject(req.Name, req.Key, projectDomain.ProjectStage(req.Stage), projectDomain.ProjectBusinessDomain(req.BusinessDomain), *createdByEmail, config.Database())

	if err != nil {
		log.Printf("Error creating project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"project": _project,
	  })
}

func RetrieveProject(c *gin.Context) {
	key := c.Param("key") // assumes route is like /projects/:key
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project key is required",
		})
		return
	}

	// Authorization
	allow, err := authorizationService.EnforceRoles(utils.GetUserRolesFromContext(c), "org", authorizationDomain.Organization, "read")

	if err != nil || !allow {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		return
	}

	// Retrieve project
	_project, err := project.RetrieveProject(key, config.Database())

	if err != nil {
		log.Printf("Error retrieving project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _project == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"project": _project,
	})
}

func RetrieveProjects(c *gin.Context) {
	allow, err := authorizationService.EnforceRoles(utils.GetUserRolesFromContext(c), "org", authorizationDomain.Organization, "read")

	if err != nil || !allow {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		return
	}

	// Retrieve projects
	_projects, err := project.RetrieveProjects(config.Database())

	if err != nil {
		log.Printf("Error retrieving projects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"projects": _projects,
	})
}

func UpdateProject(c *gin.Context) {
	key := c.Param("key") // assumes route is like /projects/:key
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project key is required",
		})
		return
	}

    var req struct {
		Name           *string                         `json:"name,omitempty"`
		Key           *string                         `json:"key,omitempty"`
		Stage          *projectDomain.ProjectStage     `json:"stage,omitempty"`
		BusinessDomain *projectDomain.ProjectBusinessDomain `json:"businessDomain,omitempty"`
		CreatedByEmail *string                         `json:"createdByEmail,omitempty"`
	}

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updatedProject, err := project.UpdateProject(
        key,
        req.Name,
		req.Key,
        req.Stage,
        req.BusinessDomain,
        config.Database(),
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updatedProject)
}