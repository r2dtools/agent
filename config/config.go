package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/unknwon/com"
)

const (
	port = 60150
)

var isDevMode = true
var Version string

type Config struct {
	LogFile,
	RootPath,
	ConfigFilePath,
	Token string
	Port      int
	IsDevMode bool
	Version   string
	vConfig   *viper.Viper
}

func GetConfig() (*Config, error) {
	var rootPath string

	if isDevMode {
		wd, err := os.Getwd()

		if err != nil {
			return nil, err
		}

		rootPath = wd

		if filepath.Base(wd) == "cmd" {
			rootPath = filepath.Dir(wd)
		}
	} else {
		executable, err := os.Executable()

		if err != nil {
			return nil, err
		}

		rootPath = filepath.Dir(executable)
	}

	vConfig := viper.New()
	vConfig.SetDefault("Port", port)

	configPath := filepath.Join(rootPath, "config")
	configFilePath := filepath.Join(configPath, "params.yaml")

	if com.IsExist(configFilePath) {
		vConfig.SetConfigType("yaml")
		vConfig.SetConfigName("params")
		vConfig.AddConfigPath(configPath)

		if err := vConfig.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	if Version == "" {
		Version = "dev"
	}

	return &Config{
		Port:           vConfig.GetInt("Port"),
		LogFile:        filepath.Join(rootPath, "sslbot.log"),
		Token:          vConfig.GetString("Token"),
		RootPath:       rootPath,
		ConfigFilePath: configFilePath,
		IsDevMode:      isDevMode,
		vConfig:        vConfig,
		Version:        Version,
	}, nil
}

func (c *Config) Merge() error {
	return c.vConfig.MergeInConfig()
}

func (c *Config) GetVarDirAbsPath() string {
	return "/usr/local/r2dtools/var"
}

func (c *Config) GetPathInsideVarDir(path ...string) string {
	parts := []string{c.GetVarDirAbsPath()}
	parts = append(parts, path...)

	return filepath.Join(parts...)
}

func (c *Config) IsSet(key string) bool {
	return c.vConfig.IsSet(key)
}

func (c *Config) GetString(key string) string {
	return c.vConfig.GetString(key)
}

func (c *Config) GetInt(key string) int {
	return c.vConfig.GetInt(key)
}

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
