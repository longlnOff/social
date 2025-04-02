package main

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Set the level encoder to uppercase or colored uppercase based on environment
	development := false
	if os.Getenv("ENVIRONMENT") == "development" {
		development = true
		encoderCfg.EncodeLevel = ColorLevelEncoder // Use colors in development
	} else {
		encoderCfg.EncodeLevel = UppercaseLevelEncoder // Just uppercase in production
	}

	// Choose console encoder for better readability with colors in development
	encoding := "json"
	if development {
		encoding = "console"
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       development,
		DisableCaller:     false,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          encoding,
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	return zap.Must(config.Build())
}

// Custom level encoder that converts log levels to uppercase
func UppercaseLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(strings.ToUpper(l.String()))
}

// Custom level encoder that adds colors to log levels
func ColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	// Color codes
	const (
		colorReset  = "\033[0m"
		colorRed    = "\033[31m"
		colorYellow = "\033[33m"
		colorBlue   = "\033[36m"
		colorGray   = "\033[37m"
	)

	// Apply color based on level
	var levelColor string
	switch l {
	case zapcore.DebugLevel:
		levelColor = colorGray
	case zapcore.InfoLevel:
		levelColor = colorBlue
	case zapcore.WarnLevel:
		levelColor = colorYellow
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		levelColor = colorRed
	default:
		levelColor = colorReset
	}

	// Format as colored uppercase level
	enc.AppendString(levelColor + strings.ToUpper(l.String()) + colorReset)
}
