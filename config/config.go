package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	defaultPort     = 60150
	defaultCaServer = "https://acme-v02.api.letsencrypt.org/directory"
	varDirPath      = "/usr/local/r2dtools/var"
)

var isDevMode = true
var Version string

type Config struct {
	LogFile   string
	Port      int
	Token     string
	IsDevMode bool
	Version   string
	rootPath  string
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

	viper.AutomaticEnv()
	viper.SetEnvPrefix("sslbot")

	if Version == "" {
		Version = "dev"
	}

	port := viper.GetInt("port")

	if port == 0 {
		port = defaultPort
	}

	tokenPath := filepath.Join(rootPath, "token")
	tokenFile, err := os.OpenFile(tokenPath, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	defer tokenFile.Close()

	token, err := io.ReadAll(tokenFile)

	if err != nil {
		return nil, err
	}

	return &Config{
		Port:      port,
		LogFile:   filepath.Join(rootPath, "sslbot.log"),
		Token:     string(token),
		rootPath:  rootPath,
		IsDevMode: isDevMode,
		Version:   Version,
	}, nil
}

func (c *Config) GetCaServer() string {
	ca := viper.GetString("ca_server")

	if ca == "" {
		ca = defaultCaServer
	}

	return ca
}

func (c *Config) GetLegoBinPath() string {
	return filepath.Join(c.rootPath, "lego")
}

func (c *Config) GetTokenPath() string {
	return filepath.Join(c.rootPath, "token")
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
