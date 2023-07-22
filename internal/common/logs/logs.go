package logs

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
)

func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			_ = l.Debug().Str("msg", msg).Msg("")
		case logging.LevelInfo:
			_ = l.Info().Str("msg", msg).Msg("")
		case logging.LevelWarn:
			_ = l.Warn().Str("msg", msg).Msg("")
		case logging.LevelError:
			_ = l.Error().Str("msg", msg).Msg("")
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
