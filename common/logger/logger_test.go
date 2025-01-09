package logger

import (
	"context"
	"testing"
)

func TestL(t *testing.T) {
	L().Info("hello world")
	L().Debug("hello world")
}

func TestCtx(t *testing.T) {
	Ctx(context.Background()).Infow("hello world", "a", 1, "b", 2)
	
}
