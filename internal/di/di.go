package di

import (
	"secret_keeper/encryptor"
	"secret_keeper/password"
	"secret_keeper/utils"

	"github.com/samber/do"
	"github.com/spf13/viper"
)

func provideConfig(path string) func(*do.Injector) (*utils.Config, error) {
	return func(i *do.Injector) (*utils.Config, error) {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}

		config := &utils.Config{}
		return config, viper.Unmarshal(config)
	}
}

func provideEncryptor(i *do.Injector) (encryptor.Encryptor, error) {
	config := do.MustInvoke[*utils.Config](i)
	return encryptor.NewSimpleEncryptor(config.SECRET_KEY, config.IV), nil
}

func provideHasher(i *do.Injector) (password.PassowrdHasher, error) {
	return password.NewSimplePasswordHasher(), nil
}

func ProvideDeps(configPath string) error {
	do.Provide(nil, provideConfig(configPath))
	do.Provide(nil, provideEncryptor)
	do.Provide(nil, provideHasher)
	return nil
}