package logging

import (
	"compartment/pkg/configuration"
	"log/slog"
)

func Error(message string) {
	logger().Error(message)
}

func Info(message string) {
	logger().Info(message)
}

func Debug(message string) {
	logger().Debug(message)
}

func logger() *slog.Logger {
	return configuration.Get().Logger
}
