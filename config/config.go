package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const prodMode = true

// Config stores agent configuration params
type Config struct {
	LogFile,
	ExecutablePath,
	Token string
	LogLevel, Port int
}

var config *Config

// GetConfig returns agent config
func GetConfig() *Config {
	return config
}

func init() {
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
	vConfig.SetConfigType("yaml")
	vConfig.SetConfigName("params")
	vConfig.AddConfigPath(filepath.Join(executablePath, "config"))
	viper.AutomaticEnv()

	if err := vConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	config = &Config{
		Port:           vConfig.GetInt("Port"),
		LogFile:        vConfig.GetString("LogFile"),
		LogLevel:       vConfig.GetInt("LogLevel"),
		Token:          vConfig.GetString("Token"),
		ExecutablePath: executablePath,
	}
}

// GetLoggerFileAbsPath returns absolute path to logger file
func (c *Config) GetLoggerFileAbsPath() string {
	return filepath.Join(c.ExecutablePath, c.LogFile)
}

// GetScriptsDirAbsPath returns absolute path to scripts directory
func (c *Config) GetScriptsDirAbsPath() string {
	return filepath.Join(c.ExecutablePath, "scripts")
}
