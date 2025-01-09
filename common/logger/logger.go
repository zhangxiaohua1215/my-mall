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

// func Ctx(ctx context.Context) Logger {
// 	var fileds []any
// 	if traceId := ctx.Value("traceid"); traceId != nil {
// 		fileds = append(fileds, "traceid", traceId)
// 	}
// 	if spanId := ctx.Value("spanid"); spanId != nil {
// 		fileds = append(fileds, "spanid", spanId)
// 	}
// 	if pSpanId := ctx.Value("pspanid"); pSpanId != nil {
// 		fileds = append(fileds, "pspanid", pSpanId)
// 	}
// 	return  defaultLogger.With(fileds...)
// }

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
