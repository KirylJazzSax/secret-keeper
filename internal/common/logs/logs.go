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
			l.Debug().Str("msg", msg).Msg("")
		case logging.LevelInfo:
			l.Info().Str("msg", msg).Msg("")
		case logging.LevelWarn:
			l.Warn().Str("msg", msg).Msg("")
		case logging.LevelError:
			l.Error().Str("msg", msg).Msg("")
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
