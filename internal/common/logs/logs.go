package logs

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			_ = log.Debug().Msg(msg)
		case logging.LevelInfo:
			_ = log.Info().Msg(msg)
		case logging.LevelWarn:
			_ = log.Warn().Msg(msg)
		case logging.LevelError:
			_ = log.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
