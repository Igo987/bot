package logger

import (
	"log/slog"
	"os"
	"time"
)

func init() {
	Log = LogStart()
}

var Log *slog.Logger

const TimeFormat = "15:04:05"
const DateFormat = "2006-01-02"

var ConfigureTime = slog.Group(
	"date",
	slog.String("", time.Now().Format(DateFormat)),
	slog.String("time", time.Now().Format(TimeFormat)),
)

func replace(groups []string, a slog.Attr) slog.Attr {
	if a.Key != slog.TimeKey || len(groups) != 0 {
		return a
	}
	return slog.Attr{}
}

type Logger struct {
	*slog.Logger
}

var defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level:       slog.LevelDebug,
	ReplaceAttr: replace,
	AddSource:   true,
}))
var logStart = func() *slog.Logger {
	l := defaultLogger

	return l.With(ConfigureTime)
}()

func LogStart() *slog.Logger {
	return logStart
}
