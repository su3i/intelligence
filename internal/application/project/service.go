package project

import (
	"errors"

	"github.com/darksuei/suei-intelligence/internal/application/account"
	"github.com/darksuei/suei-intelligence/internal/config"
	"github.com/darksuei/suei-intelligence/internal/domain/project"
	projectDomain "github.com/darksuei/suei-intelligence/internal/domain/project"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/database"
)

func NewProject(name string, key string,stage projectDomain.ProjectStage, businessDomain projectDomain.ProjectBusinessDomain, createdByEmail string, cfg *config.DatabaseConfig) (*project.Project, error) {
	_projectRepository := database.NewProjectRepository(cfg)

	_project, err := _projectRepository.FindOneByKey(key)

	if err != nil {
		return nil, err
	}

	if _project != nil {
		return nil, errors.New("Project already exists.")
	}

	createdByAccount, err := account.RetrieveAccount(createdByEmail, cfg)

	if err != nil {
		return nil, errors.New("Failed to get account")
	}

	createdBy := map[string]string{
		"Email": createdByEmail,
		"Name": createdByAccount.Name,
	}

	_project = &project.Project{
		Name:  name,
		Key: key,
		Status: project.Active,
		Stage:   stage,
		BusinessDomain: businessDomain,
		CreatedBy: createdBy,
	}

	return _projectRepository.Create(_project)
}

func RetrieveProject(key string, cfg *config.DatabaseConfig) (*project.Project, error) {
	_projectRepository := database.NewProjectRepository(cfg)

	return _projectRepository.FindOneByKey(key)
}

func RetrieveProjects(cfg *config.DatabaseConfig) (*[]project.Project, error) {
	_projectRepository := database.NewProjectRepository(cfg)

	return _projectRepository.Find()
}

func UpdateProject(
    key string,
    name *string,
    newKey *string,
    stage *projectDomain.ProjectStage,
    businessDomain *projectDomain.ProjectBusinessDomain,
    cfg *config.DatabaseConfig,
) (*project.Project, error) {
    _projectRepository := database.NewProjectRepository(cfg)

    // Find existing project
    _project, err := _projectRepository.FindOneByKey(key)
    if err != nil {
        return nil, err
    }

    if _project == nil {
        return nil, errors.New("Project not found")
    }

    // Update fields only if provided
    if name != nil {
        _project.Name = *name
    }
	if newKey != nil {
        _project.Key = *newKey
    }
    if stage != nil {
        _project.Stage = *stage
    }
    if businessDomain != nil {
        _project.BusinessDomain = *businessDomain
    }

    // Save updated project
    if err := _projectRepository.Update(_project); err != nil {
        return nil, err
    }

    return _project, nil
}