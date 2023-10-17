package log

import (
	"context"
	"time"

	"github.com/ByPikod/go-crypto/tree/crypto/core"
	zaploki "github.com/paul-milne/zap-loki"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

// Init logger
func InitializeLogger() {
	var err error
	logger, err = CreateLogger()
	if err != nil {
		panic(err)
	}
	Info("Loki connection established")
}

// Create logger instance
func CreateLogger() (*zap.Logger, error) {
	//  Create "Loki" hook
	config := zap.NewProductionConfig()
	lokiInfo := core.Config.Loki
	loki := zaploki.New(
		context.Background(),
		zaploki.Config{
			Url:          lokiInfo.Host,
			BatchMaxSize: 1000,
			BatchMaxWait: 10 * time.Second,
			Labels:       map[string]string{"app": "Crypto"},
		},
	)

	// Create logger
	logger, err := loki.WithCreateLogger(config)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// Logs the string: An error ocurred at %v controller.
func ControllerError(controller string, err error) {
	Error("An error ocurred at a controller", zap.String("controller", controller), zap.Error(err))
}

// Logs the text and puts "zap.Error" field.
func QuickError(text string, err error) {
	Error(text, zap.Error(err))
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}
