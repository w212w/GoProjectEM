// internal/logger/logger.go

package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

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
