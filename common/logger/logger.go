package logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger = *zap.SugaredLogger
	
var ctxKeys []string

func RegisterCtxKeys(keys ...string) {
	ctxKeys = append(ctxKeys, keys...)
}


func Ctx(ctx context.Context) Logger {
	var fileds []any
	for _, ctxkey := range ctxKeys {
		if v := ctx.Value(ctxkey); v != nil {
			fileds = append(fileds, ctxkey, v)
		}
	}
	return  defaultLogger.With(fileds...)
}

func L() Logger {
	return defaultLogger
}
