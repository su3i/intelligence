package account

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model

	FullName         string `gorm:"unique;not null"`
	Email         string `gorm:"unique;not null"`
	PasswordEnc		string
	Role 		AccountRole `gorm:"type:text;not null"`
}