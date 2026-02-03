package datasource

import (
	"gorm.io/gorm"
)

type DataSource struct {
	gorm.Model

	ProjectID uint   `gorm:"not null;index"` // <- foreign key to Project
}