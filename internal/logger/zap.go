package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: logger}
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.logger.Error(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	z.logger.Fatal(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
