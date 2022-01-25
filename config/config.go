package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/unknwon/com"
)

const (
	prodMode = false
	port     = 60150
	logFile  = "var/log/r2dtools.log"
	logLevel = 4
)

// Config stores agent configuration params
type Config struct {
	LogFile,
	ExecutablePath,
	ConfigPath,
	Token string
	LogLevel, Port int
	vConfig        *viper.Viper
}

var config *Config

// GetConfig returns agent config
func GetConfig() *Config {
	if config != nil {
		return config
	}

	var executablePath string
	var err error

	if prodMode {
		executable, err := os.Executable()

		if err != nil {
			panic(err)
		}

		executablePath = filepath.Dir(executable)
	} else {
		executablePath, err = os.Getwd()

		if err != nil {
			panic(err)
		}
	}

	vConfig := viper.New()
	vConfig.SetDefault("Port", port)
	vConfig.SetDefault("LogFile", logFile)
	vConfig.SetDefault("LogLevel", logLevel)

	configPath := filepath.Join(executablePath, "config")
	configFilePath := filepath.Join(configPath, "params.yaml")

	if com.IsExist(configFilePath) {
		vConfig.SetConfigType("yaml")
		vConfig.SetConfigName("params")
		vConfig.AddConfigPath(configPath)
		viper.AutomaticEnv()

		if err := vConfig.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	config = &Config{
		Port:           vConfig.GetInt("Port"),
		LogFile:        vConfig.GetString("LogFile"),
		LogLevel:       vConfig.GetInt("LogLevel"),
		Token:          vConfig.GetString("Token"),
		ExecutablePath: executablePath,
		ConfigPath:     configPath,
		vConfig:        vConfig,
	}

	return config
}

// GetLoggerFileAbsPath returns absolute path to logger file
func (c *Config) GetLoggerFileAbsPath() string {
	return filepath.Join(c.ExecutablePath, c.LogFile)
}

// Merge merges config already loaded in memory with existing one
func (c *Config) Merge() error {
	return c.vConfig.MergeInConfig()
}

// GetVarDirAbsPath returns absolute path to var directory
func (c *Config) GetVarDirAbsPath() string {
	return filepath.Join(c.ExecutablePath, "var")
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
