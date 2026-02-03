package appconfig

import (
	"gorm.io/gorm"
)

type AppConfig struct {
	gorm.Model

	BootstrapToken         string `gorm:"unique;not null"`
}