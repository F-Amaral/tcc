package log

import "go.uber.org/zap"

type Field = zap.Field

func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

func Err(err error) Field {
	return zap.Error(err)
}
