package log

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/F-Amaral/tcc/internal/apierrors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

type wrapper struct {
	level  zapcore.Level
	logger Logger
}

func Debug(ctx context.Context) wrapper {
	return wrapper{
		level:  zap.DebugLevel,
		logger: getLogger(ctx),
	}
}

func Warn(ctx context.Context) wrapper {
	return wrapper{
		level:  zap.WarnLevel,
		logger: getLogger(ctx),
	}
}

func Error(ctx context.Context) wrapper {
	return wrapper{
		level:  zap.ErrorLevel,
		logger: getLogger(ctx),
	}
}

func Info(ctx context.Context) wrapper {
	return wrapper{
		level:  zap.InfoLevel,
		logger: getLogger(ctx),
	}
}

func (l wrapper) zap(msg string, fields ...Field) {
	switch l.level {
	case zap.ErrorLevel:
		l.logger.Errorln(msg, fields)
	case zap.InfoLevel:
		l.logger.Infoln(msg, fields)
	case zap.WarnLevel:
		l.logger.Warnln(msg, fields)
	default:
		l.logger.Debugln(msg, fields)
	}
}

// LogApiError Put important root cause information into the apiError
func (l wrapper) LogApiError(err apierrors.ApiError) {
	l.zap(fmt.Sprintf("ApiError: %+v", err))
}

func (l wrapper) Log(msg string, fields ...Field) {
	l.zap(msg, fields...)
}

func (l wrapper) LogError(err error) {
	l.zap(err.Error())
}

func (l wrapper) Json(val any, fields ...Field) {
	b, _ := json.Marshal(val)
	l.zap(string(b), fields...)
}
func getLogger(ctx context.Context) Logger {
	logger := FromContext(ctx)
	if logger != nil {
		return logger
	}

	return zap.NewNop().Sugar()
}
