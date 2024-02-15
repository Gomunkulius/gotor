package internal

import (
	"gotor/internal/log"
	"log/slog"
)

var (
	Logger *slog.Logger
	Port   int = 6484
)

func InitGlobal() {
	Logger = log.GetLogger()
	Logger.Info("Successfully init global variables")
}
