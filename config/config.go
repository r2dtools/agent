package config

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	defaultPort       = 60150
	defaultCaServer   = "https://acme-v02.api.letsencrypt.org/directory"
	defaultVarDirPath = "/usr/local/r2dtools/sslbot/var"
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
	VarDirPath     string
	rootPath       string
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

	configFilePath := filepath.Join(rootPath, "config.yaml")
	configFile, err := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	viper.AddConfigPath(filepath.Dir(configFilePath))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("sslbot")

	viper.SetDefault("port", defaultPort)
	viper.SetDefault("ca_server", defaultCaServer)
	viper.SetDefault("var_dir_path", defaultVarDirPath)

	if err := viper.ReadConfig(configFile); err != nil {
		panic(err)
	}

	if Version == "" {
		Version = "dev"
	}

	if isDevMode {
		viper.Set("var_dir_path", filepath.Join(rootPath, "var"))
	}

	config := &Config{
		LogFile:        filepath.Join(rootPath, "sslbot.log"),
		LegoBinPath:    filepath.Join(rootPath, "lego"),
		ConfigFilePath: configFilePath,
		rootPath:       rootPath,
		IsDevMode:      isDevMode,
		Version:        Version,
	}
	setDynamicParams(config)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		setDynamicParams(config)
	})

	return config, nil
}

func (c *Config) GetPathInsideVarDir(path ...string) string {
	parts := []string{c.VarDirPath}
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

func setDynamicParams(c *Config) {
	c.Port = viper.GetInt("port")
	c.Token = viper.GetString("token")
	c.CaServer = viper.GetString("ca_server")
	c.VarDirPath = viper.GetString("var_dir_path")
}
