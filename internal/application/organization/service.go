package organization

import (
	"errors"

	"github.com/darksuei/suei-intelligence/internal/config"
	"github.com/darksuei/suei-intelligence/internal/domain/organization"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/database"
)

func NewOrganization(name string, key string, scope organization.OrgScope, cfg *config.DatabaseConfig) (*organization.Organization, error) {
	_organizationRepository := database.NewOrganizationRepository(cfg)

	_organization, err := _organizationRepository.FindOne(key)

	if err != nil {
		return nil, err
	}

	if _organization != nil {
		return nil, errors.New("Organization already exists.")
	}

	_organization = &organization.Organization{
		Name:  name,
		Key:   key,
		Scope: scope,
	}

	return _organizationRepository.Create(_organization)
}

func RetrieveOrganization(key string, cfg *config.DatabaseConfig) (*organization.Organization, error) {
	_organizationRepository := database.NewOrganizationRepository(cfg)

	return _organizationRepository.FindOne(key)
}