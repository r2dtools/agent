package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/unknwon/com"
)

const (
	devMode = "development"
	port    = 60150
)

// Config stores agent configuration params
type Config struct {
	LogFile,
	ExecutablePath,
	ConfigPath,
	Token string
	Port      int
	IsDevMode bool
	vConfig   *viper.Viper
}

// GetConfig returns agent config
func GetConfig() (*Config, error) {
	env := os.Getenv("R2DTOOLS_AGENT_MODE")
	isDevMode := env == devMode

	var executablePath string
	var err error

	if isDevMode {
		executablePath, err = os.Getwd()

		if err != nil {
			return nil, err
		}
	} else {
		executable, err := os.Executable()

		if err != nil {
			return nil, err
		}

		executablePath = filepath.Dir(executable)

	}

	vConfig := viper.New()
	vConfig.SetDefault("Port", port)

	configPath := filepath.Join(executablePath, "config")
	configFilePath := filepath.Join(configPath, "params.yaml")

	if com.IsExist(configFilePath) {
		vConfig.SetConfigType("yaml")
		vConfig.SetConfigName("params")
		vConfig.AddConfigPath(configPath)

		if err := vConfig.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	return &Config{
		Port:           vConfig.GetInt("Port"),
		LogFile:        filepath.Join(executablePath, "r2dtools.log"),
		Token:          vConfig.GetString("Token"),
		ExecutablePath: executablePath,
		ConfigPath:     configPath,
		IsDevMode:      isDevMode,
		vConfig:        vConfig,
	}, nil
}

// Merge merges config already loaded in memory with existing one
func (c *Config) Merge() error {
	return c.vConfig.MergeInConfig()
}

// GetVarDirAbsPath returns absolute path to var directory
func (c *Config) GetVarDirAbsPath() string {
	return "/usr/local/r2dtools/var"
}

// IsSet checks if key exists in config
func (c *Config) IsSet(key string) bool {
	return c.vConfig.IsSet(key)
}

// GetString returns string value by key
func (c *Config) GetString(key string) string {
	return c.vConfig.GetString(key)
}

// GetInt returns int value by key
func (c *Config) GetInt(key string) int {
	return c.vConfig.GetInt(key)
}

// GetModuleVarDir returns module var directiry absolute path
func (c *Config) GetModuleVarAbsDir(id string) string {
	return filepath.Join(c.GetVarDirAbsPath(), "modules", fmt.Sprintf("%s-module", id))
}

// ToMap returns all settings as map
func (c *Config) ToMap() map[string]string {
	settings := c.vConfig.AllSettings()
	options := make(map[string]string)

	for key, value := range settings {
		if strValue, ok := value.(string); ok {
			options[key] = strValue
		}
	}

	return options
}
