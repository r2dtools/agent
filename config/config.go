package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	defaultPort     = 60150
	defaultCaServer = "https://acme-v02.api.letsencrypt.org/directory"
	varDirPath      = "/usr/local/r2dtools/sslbot/var"
)

var isDevMode = true
var Version string

type Config struct {
	LogFile        string
	Port           int
	Token          string
	IsDevMode      bool
	Version        string
	LegoBinPath    string
	CaServer       string
	ConfigFilePath string
	rootPath       string
	viper          *viper.Viper
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
	configFilePath := filepath.Join(rootPath, "config.yaml")
	configFile, err := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	vConfig.AddConfigPath(configFilePath)
	vConfig.SetConfigType("yaml")
	vConfig.SetConfigName("params")

	if err := vConfig.ReadConfig(configFile); err != nil {
		panic(err)
	}

	vConfig.AutomaticEnv()
	vConfig.SetEnvPrefix("sslbot")

	vConfig.SetDefault("port", defaultPort)
	vConfig.SetDefault("ca_server", defaultCaServer)

	if Version == "" {
		Version = "dev"
	}

	return &Config{
		Port:           vConfig.GetInt("port"),
		LogFile:        filepath.Join(rootPath, "sslbot.log"),
		Token:          vConfig.GetString("token"),
		LegoBinPath:    filepath.Join(rootPath, "lego"),
		CaServer:       vConfig.GetString("ca_server"),
		ConfigFilePath: configFilePath,
		rootPath:       rootPath,
		IsDevMode:      isDevMode,
		Version:        Version,
		viper:          vConfig,
	}, nil
}

func (c *Config) GetPathInsideVarDir(path ...string) string {
	parts := []string{varDirPath}
	parts = append(parts, path...)

	return filepath.Join(parts...)
}

func (c *Config) ToMap() map[string]string {
	settings := viper.AllSettings()
	options := make(map[string]string)

	for key, value := range settings {
		if strValue, ok := value.(string); ok {
			options[key] = strValue
		}
	}

	return options
}
