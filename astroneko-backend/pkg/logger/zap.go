package logger

import "go.uber.org/zap"

func NewZapLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
