package logger

import (
	"log"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func LogInit(modeLog string) {
	var handler slog.Handler

	switch modeLog {
	case "debug":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "jsonDebug":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "jsonInfo":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		log.Fatal("not init modelog: ", modeLog)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
	return
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key: "error",
		Value: slog.StringValue(err.Error()),
	}
}