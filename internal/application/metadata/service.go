package metadata

import (
	"log"

	"github.com/darksuei/suei-intelligence/internal/config"
	"github.com/darksuei/suei-intelligence/internal/domain/metadata"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/database"
)

// Check if bootstrap token exists in database
// If it does not, create it
// If it does, validate it
// If validation fails, throw error
func LoadBootstrapToken(bootstrapToken string, cfg *config.DatabaseConfig) error {
	log.Print("Loading bootstrap token..")

	_metadataRepository := database.NewMetadataRepository(cfg)

	_metadata, err := _metadataRepository.FindOne();

	if err != nil {
		return err
	}

	if _metadata != nil {
		if _metadata.BootstrapToken == bootstrapToken {
			return nil
		}
		panic("Do not change bootstrap token!")
	}

	var _newMetadata metadata.Metadata

	_newMetadata.BootstrapToken = bootstrapToken

	_metadataRepository.Create(&_newMetadata)

	return nil
}

func SetLanguage(language string, cfg *config.DatabaseConfig) error {
	_metadataRepository := database.NewMetadataRepository(cfg)

	_metadata, err := _metadataRepository.FindOne();

	if err != nil {
		return err
	}

	if _metadata != nil {
		_metadata.Language = language
	}

	return _metadataRepository.Update(_metadata)
}