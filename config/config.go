package config

import "github.com/spf13/viper"

// Config stores agent configuration params
type Config struct {
	LogFile,
	Token string
	LogLevel, Port int
}

var config *Config

// GetConfig returns agent config
func GetConfig() *Config {
	return config
}

func init() {
	vConfig := viper.New()
	vConfig.SetConfigType("yaml")
	vConfig.SetConfigName("params")
	vConfig.AddConfigPath("config/")
	viper.AutomaticEnv()

	if err := vConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	config = &Config{
		Port:     vConfig.GetInt("Port"),
		LogFile:  vConfig.GetString("LogFile"),
		LogLevel: vConfig.GetInt("LogLevel"),
		Token:    vConfig.GetString("Token"),
	}
}
