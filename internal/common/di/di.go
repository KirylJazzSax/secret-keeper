package di

import (
	"secret_keeper/encryptor"
	"secret_keeper/gapi"
	"secret_keeper/internal/common/utils"
	"secret_keeper/password"
	"secret_keeper/repository"
	"secret_keeper/token"
	"secret_keeper/validation"

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
	return encryptor.NewSimpleEncryptor(config.SecretKey, config.IV), nil
}

func provideHasher(i *do.Injector) (password.PassowrdHasher, error) {
	return password.NewSimplePasswordHasher(), nil
}

func provideMaker(i *do.Injector) (token.Maker, error) {
	config := do.MustInvoke[*utils.Config](i)
	return token.NewPasetoMaker(config.SymmetricKey)
}

func provideValidator(i *do.Injector) (validation.Validator, error) {
	return validation.NewSimpleValidator(), nil
}

func provideRepository(i *do.Injector) (repository.Repository, error) {
	config := do.MustInvoke[*utils.Config](i)
	return repository.NewBoltRepository(config.DbUrl)
}

func provideServer(i *do.Injector) (*gapi.Server, error) {
	tokenManager := do.MustInvoke[token.Maker](i)
	repo := do.MustInvoke[repository.Repository](i)
	config := do.MustInvoke[*utils.Config](i)
	return gapi.NewServer(repo, tokenManager, config), nil
}

func ProvideDeps(configPath string) error {
	do.Provide(nil, provideConfig(configPath))
	do.Provide(nil, provideEncryptor)
	do.Provide(nil, provideHasher)
	do.Provide(nil, provideMaker)
	do.Provide(nil, provideValidator)
	do.Provide(nil, provideRepository)
	do.Provide(nil, provideServer)
	return nil
}
