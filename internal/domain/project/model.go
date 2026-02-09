package project

import (
	"gorm.io/gorm"

	"github.com/darksuei/suei-intelligence/internal/domain/datasource"
)

type Project struct {
	gorm.Model

	Name            string                   `gorm:"unique;not null"`
	Key             string                   `gorm:"unique;not null"`
	Status          ProjectStatus            `gorm:"type:text;not null"`
	Stage           ProjectStage             `gorm:"type:text;not null"`
	BusinessDomain  ProjectBusinessDomain    `gorm:"type:text;not null"`
	CreatedBy       map[string]string 		 `gorm:"type:jsonb;serializer:json;default:'{}'"`
	DataSources     []datasource.DataSource  `gorm:"foreignKey:ProjectID"`
}
