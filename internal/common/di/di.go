package di

import (
	"os"
	"time"

	"github.com/KirylJazzSax/secret-keeper/internal/common/encryptor"
	"github.com/KirylJazzSax/secret-keeper/internal/common/password"
	"github.com/KirylJazzSax/secret-keeper/internal/common/token"
	"github.com/KirylJazzSax/secret-keeper/internal/common/utils"
	"github.com/KirylJazzSax/secret-keeper/internal/common/validation"

	"github.com/samber/do"
	"github.com/spf13/viper"
)

func provideEnvConfig(i *do.Injector) (*utils.Config, error) {
	dur, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return nil, err
	}
	return &utils.Config{
		HttpPort:            os.Getenv("HTTP_PORT"),
		SymmetricKey:        os.Getenv("SYMMETRIC_KEY"),
		SecretKey:           os.Getenv("SECRET_KEY"),
		IV:                  os.Getenv("IV"),
		AccessTokenDuration: dur,
		SrvType:             os.Getenv("SERVER_TYPE"),
		DbUsername:          os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		DbPassword:          os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
		DbPort:              os.Getenv("DB_PORT"),
		GrpcEndpoint:        os.Getenv("GRPC_ENDPOINT"),
	}, nil
}

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

// func provideRepository(i *do.Injector) (repository.Repository, error) {
// 	config := do.MustInvoke[*utils.Config](i)
// 	return repository.NewBoltRepository(config.DbUrl)
// }

// func provideServer(i *do.Injector) (*gapi.Server, error) {
// 	tokenManager := do.MustInvoke[token.Maker](i)
// 	repo := do.MustInvoke[repository.Repository](i)
// 	config := do.MustInvoke[*utils.Config](i)
// 	return gapi.NewServer(repo, tokenManager, config), nil
// }

func ProvideDeps() error {
	do.Provide(nil, provideEnvConfig)
	do.Provide(nil, provideEncryptor)
	do.Provide(nil, provideHasher)
	do.Provide(nil, provideMaker)
	do.Provide(nil, provideValidator)
	return nil
}
