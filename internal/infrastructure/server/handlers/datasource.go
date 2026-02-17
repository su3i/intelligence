package handlers

import (
	"log"
	"net/http"
	"strconv"

	authorizationService "github.com/darksuei/suei-intelligence/internal/application/authorization"
	datasourceService "github.com/darksuei/suei-intelligence/internal/application/datasource"
	"github.com/darksuei/suei-intelligence/internal/config"
	authorizationDomain "github.com/darksuei/suei-intelligence/internal/domain/authorization"
	datasourceDomain "github.com/darksuei/suei-intelligence/internal/domain/datasource"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/etl"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SupportedDatasources(c *gin.Context) {
    stripped := make([]map[string]interface{}, len(datasourceDomain.SupportedDatasources))
    for i, ds := range datasourceDomain.SupportedDatasources {
        copy := make(map[string]interface{})
        for k, v := range ds {
            if k != "form" {
                copy[k] = v
            }
        }
        stripped[i] = copy
    }

    c.JSON(http.StatusOK, gin.H{
        "message":     "success",
        "datasources": stripped,
    })
}

func SupportedDatasource(c *gin.Context) {
	sourceType := c.Param("sourceType")

	for _, ds := range datasourceDomain.SupportedDatasources {
		if ds["sourceType"] == sourceType {
			c.JSON(http.StatusOK, gin.H{
				"message":    "success",
				"datasource": ds,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "datasource not found",
	})
}

func NewDatasource(c *gin.Context) {
	var req struct {
		SourceType string `json:"sourceType" binding:"required"`
		Configuration map[string]interface{} `json:"configuration" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Validation failed.",
			"errors": utils.FormatValidationErrors(err),
		})
		return
	}

	projectKey := c.Param("key") // assumes route is like /projects/:key
	if projectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project key is required",
		})
		return
	}

	errs, err := datasourceDomain.ValidateInput(req.SourceType, req.Configuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Validation failed.",
			"errors":  errs,
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

	configuration := req.Configuration
	configuration["sourceType"] = req.SourceType

	sourceId, err := etl.GetInstance().CreateSourceConnection(uuid.New().String(), configuration)

	if err != nil {
		log.Printf("Error creating datasource: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Test connection
	err = etl.GetInstance().TestSourceConnection(*sourceId)

	// If connection fails - delete ETL source
	if err != nil {
		// Rollback CREATED ETL source
		etl.GetInstance().DeleteSourceConnection(*sourceId)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to datasource, please check your connection details and try again.",
		})
		return
	}

	// If success, create datasource
	createdByEmail, err := utils.GetUserEmailFromContext(c)

	if err != nil || createdByEmail == nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get account",
		})
		return
	}

	// Create datasource
	_datasource, err := datasourceService.NewDatasource(projectKey, req.SourceType, *sourceId, *createdByEmail, config.Database())

	if err != nil {
		log.Printf("Error creating datasource: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"datasource": _datasource,
	})
	return
}

func RetrieveDatasources(c *gin.Context) {
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

	// Retrieve datasources
	_datasources, err := datasourceService.RetrieveDatasources(key, config.Database())

	if err != nil {
		log.Printf("Error retrieving datasources: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"datasources": _datasources,
	})
}

func DeleteDatasource(c *gin.Context) {
	key := c.Param("key") // assumes route is like /projects/:key
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project key is required",
		})
		return
	}

	idParam := c.Param("id") // /projects/:id
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project id is required",
		})
		return
	}

	datasourceID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid datasource id",
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

	// Delete datasource
	err = datasourceService.SoftDeleteDatasource(uint(datasourceID), key, config.Database())

	if err != nil {
		log.Printf("Error deleting datasource: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}