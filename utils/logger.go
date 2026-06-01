package utils

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

var wg sync.WaitGroup

func NewLogger(ctx context.Context, loggerType string) (*zap.Logger, error) {
	if logger, err := getLoggerConstructor(loggerType)(); err != nil {
		return nil, err
	} else {
		wg.Add(1)

		go func() {
			defer wg.Done()
			<-ctx.Done()

			logger.Sync()
		}()

		return logger, nil
	}
}

func getLoggerConstructor(loggerType string) func(...zap.Option) (*zap.Logger, error) {
	switch loggerType {
	case "development":
		return zap.NewDevelopment
	default:
		return zap.NewProduction
	}
}

func Wait() {
	wg.Wait()
}
