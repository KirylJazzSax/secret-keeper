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
	Port                string        `mapstructure:"PORT"`
	IV                  string        `mapstructure:"IV"`
	SecretKey           string        `mapstructure:"SECRET_KEY"`
	DbUrl               string        `mapstructure:"DB_URL"`
	HttpPort            string        `mapstructure:"HTTP_PORT"`
	Cors                string        `mapstructure:"CORS"`
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
