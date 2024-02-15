package internal

import (
	"gotor/internal/log"
	"log/slog"
)

var (
	Logger *slog.Logger
)

func InitGlobal() {
	Logger = log.GetLogger()
	Logger.Info("Successfully init global variables")
}
