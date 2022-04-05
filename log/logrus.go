package log

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

func NewLogrus(ctx context.Context) *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: false,
		DisableColors: false,
	})
	return logger
}
