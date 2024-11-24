// internal/logger/logger.go

package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log - это глобальный экземпляр логера
var Log = logrus.New()

// SetupLogger настраивает логгер с уровнем логирования из переменной окружения.
func SetupLogger() {
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "info"
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		Log.Fatalf("Неверный уровень логирования: %v", err)
	}
	Log.SetLevel(level)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Log.Debug("Logrus успешно инициализирован")
}
