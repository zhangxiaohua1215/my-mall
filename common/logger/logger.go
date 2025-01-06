package logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger struct {
	base *zap.SugaredLogger
}

func Ctx(ctx context.Context) *Logger {
	var traceId, spanId, pSpanId string
	if ctx.Value("traceid") != nil {
		traceId = ctx.Value("traceid").(string)
	}
	if ctx.Value("spanid") != nil {
		spanId = ctx.Value("spanid").(string)
	}
	if ctx.Value("psapnid") != nil {
		pSpanId = ctx.Value("pspanid").(string)
	}

	return &Logger{
		base: defaultLogger.With("traceid", traceId,
			"spanid", spanId,
			"pspanid", pSpanId),
	}
}
