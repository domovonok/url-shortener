package logger

import "context"

type RedisLogger struct {
	log Logger
}

func NewRedisLogger(log Logger) *RedisLogger {
	return &RedisLogger{log: log}
}

func (l *RedisLogger) Printf(_ context.Context, format string, v ...interface{}) {
	l.log.Info(format, Any("args", v))
}
