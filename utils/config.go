package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	GRPCServerAddress   string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	SymmetricKey        string        `mapstructure:"SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	PORT                string        `mapstructure:"PORT"`
	IV                  string        `mapstructure:"IV"`
	SECRET_KEY          string        `mapstructure:"SECRET_KEY"`
	DB_URL              string        `mapstructure:"DB_URL"`
	HTTP_PORT           string        `mapstructure:"HTTP_PORT"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
