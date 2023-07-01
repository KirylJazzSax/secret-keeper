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
	HttpPort            string        `mapstructure:"HTTP_PORT"`
	Cors                string        `mapstructure:"CORS"`
	SrvType             string        `mapstructure:"SERVER_TYPE"`
	DbUsername          string        `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	DbPassword          string        `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`
	DbPort              string        `mapstructure:"DB_PORT"`
	GrpcEndpoint        string        `mapstructure:"GRPC_ENDPOINT"`
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
