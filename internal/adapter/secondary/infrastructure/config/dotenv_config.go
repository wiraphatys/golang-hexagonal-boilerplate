package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"basedir/internal/core/shared/config"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Env           string
	Name          string
	Host          string
	Port          string
	BaseApiPrefix string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

type AppEnvConfig struct {
	serverCfg *ServerConfig
	dbCfg     *DatabaseConfig
}

var _ config.ServerConfigProvider = (*AppEnvConfig)(nil)
var _ config.DatabaseConfigProvider = (*AppEnvConfig)(nil)
var _ config.AppConfigProvider = (*AppEnvConfig)(nil)

var (
	loadedConfig *AppEnvConfig
	loadOnce     sync.Once
)

func LoadConfig(envFilePath ...string) (*AppEnvConfig, error) {
	var loadErr error
	loadOnce.Do(func() {
		if len(envFilePath) > 0 && envFilePath[0] != "" {
			err := godotenv.Load(envFilePath[0])
			if err != nil {
				log.Printf("INFO: .env file not found or failed to load from %s: %v. Proceeding with system environment variables and/or defaults.", envFilePath[0], err)
			}
		}

		serverCfg := &ServerConfig{
			Env:           getEnv("SERVER_ENV", "development"),
			Name:          getEnv("SERVER_NAME", "MyApp"),
			Host:          getEnv("SERVER_HOST", "0.0.0.0"),
			Port:          getEnv("SERVER_PORT", "8080"),
			BaseApiPrefix: getEnv("SERVER_BASEAPIPREFIX", "/api/v1"),
		}

		dbCfg := &DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "1234"),
			Name:     getEnv("DB_NAME", "mydatabase"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Bangkok"),
		}

		if dbCfg.User == "" {
			loadErr = fmt.Errorf("DB_USER cannot be empty")
			return
		}

		loadedConfig = &AppEnvConfig{
			serverCfg: serverCfg,
			dbCfg:     dbCfg,
		}
		log.Println("INFO: Application configuration loaded successfully.")
	})

	if loadErr != nil {
		return nil, loadErr
	}
	if loadedConfig == nil && loadErr == nil {
		return nil, fmt.Errorf("configuration was not loaded")
	}
	return loadedConfig, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func (c *AppEnvConfig) GetServerEnv() string {
	return c.serverCfg.Env
}
func (c *AppEnvConfig) GetServerName() string {
	return c.serverCfg.Name
}
func (c *AppEnvConfig) GetServerHost() string {
	return c.serverCfg.Host
}
func (c *AppEnvConfig) GetServerPort() string {
	return c.serverCfg.Port
}
func (c *AppEnvConfig) GetServerBaseApiPrefix() string {
	return c.serverCfg.BaseApiPrefix
}

func (c *AppEnvConfig) GetDBHost() string {
	return c.dbCfg.Host
}
func (c *AppEnvConfig) GetDBPort() string {
	return c.dbCfg.Port
}
func (c *AppEnvConfig) GetDBUser() string {
	return c.dbCfg.User
}
func (c *AppEnvConfig) GetDBPassword() string {
	return c.dbCfg.Password
}
func (c *AppEnvConfig) GetDBName() string {
	return c.dbCfg.Name
}
func (c *AppEnvConfig) GetDBSSLMode() string {
	return c.dbCfg.SSLMode
}
func (c *AppEnvConfig) GetDBTimezone() string {
	return c.dbCfg.Timezone
}

func (c *AppEnvConfig) GetDBDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.dbCfg.Host,
		c.dbCfg.Port,
		c.dbCfg.User,
		c.dbCfg.Password,
		c.dbCfg.Name,
		c.dbCfg.SSLMode,
		c.dbCfg.Timezone,
	)
}
