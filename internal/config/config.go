package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var (
	airbyte  *AirbyteConfig
	cache	 *CacheConfig
	casbin   *CasbinConfig
    common   *CommonConfig
    database *DatabaseConfig
)

func Initialize() {
	airbyte = &AirbyteConfig{}
    if err := envconfig.Process("", airbyte); err != nil {
		log.Fatalf("airbyte config: %v", err)
    }
	cache = &CacheConfig{}
    if err := envconfig.Process("", cache); err != nil {
		log.Fatalf("cache config: %v", err)
    }
	casbin = &CasbinConfig{}
    if err := envconfig.Process("", casbin); err != nil {
		log.Fatalf("casbin config: %v", err)
    }
	common = &CommonConfig{}
    if err := envconfig.Process("", common); err != nil {
		log.Fatalf("common config: %v", err)
    }
	database = &DatabaseConfig{}
	if err := envconfig.Process("", database); err != nil {
		log.Fatalf("database config: %v", err)
	}
}

func Airbyte() *AirbyteConfig     { return airbyte }
func Cache() *CacheConfig     { return cache }
func Casbin() *CasbinConfig     { return casbin }
func Common() *CommonConfig     { return common }
func Database() *DatabaseConfig { return database }