package logging

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

// Debug logs a message at level Debug on the standard logger.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info logs a message at level Info on the standard logger.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error logs a message at level Error on the standard logger.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Panic logs a message at level Panic on the standard logger.
// then call panic()
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Fatal logs a message at level Fatal on the standard logger
// then call os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
