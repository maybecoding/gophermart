package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// lg Общая переменная для логирования будет доступна всему коду
// Не лучшее решение, но самое простое
var lg *zerolog.Logger

// Init - Инициализация логгера
func Init(level string) {
	// Пока используем консольный вывод
	zl := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	switch level {
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zl.Debug().Msg("passed wrong error level")
	}
	zl.Debug().Str("log level", level).Msg("log initialized")

	lg = &zl
}

// Fatal starts a new message with fatal level. The os.Exit(1) function is called by the Msg method, which terminates the program immediately.
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return lg.Fatal()
}

// Error - Error starts a new message with error level.
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return lg.Error()
}

// Info - Info starts a new message with info level.
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return lg.Info()
}

// Debug - Debug starts a new message with debug level.
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return lg.Debug()
}
