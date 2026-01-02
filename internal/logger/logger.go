package logger

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Sync() error
}

type Field struct {
	Key   string
	Value any
}

func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}
