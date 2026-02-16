package etl

import (
	"sync"

	"github.com/darksuei/suei-intelligence/internal/config"
	domain "github.com/darksuei/suei-intelligence/internal/domain/etl"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/etl/airbyte"
)

var (
	instance domain.ETL
	once     sync.Once
)

// GetInstance returns a singleton ETL instance
func GetInstance() domain.ETL {
	return airbyte.Initialize(config.Airbyte())
}