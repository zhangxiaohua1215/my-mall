package logger

import (
	"my-mall/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultLogger Logger
func init() {
	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	fileWriteSyncer := getFileLogWriter()
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	var core zapcore.Core
	switch config.App.Env {
	case "test", "prod":
		// 测试环境和生产环境的日志输出到文件中
		core = zapcore.NewCore(jsonEncoder, fileWriteSyncer, zapcore.InfoLevel)
	case "dev":
		// 开发环境同时向控制台和文件输出日志， Debug级别的日志也会被输出
		core = zapcore.NewTee(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			zapcore.NewCore(jsonEncoder, fileWriteSyncer, zapcore.DebugLevel))
	}
	defaultLogger = zap.New(core, zap.AddCaller()).Sugar()
}

func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
	// 使用 lumberjack 实现 logger rotate
	lumberJackLogger := &lumberjack.Logger{
		Filename:  config.App.Log.FilePath,
		MaxSize:   config.App.Log.FileMaxSize,      // 文件最大 100 M
		MaxAge:    config.App.Log.BackUpFileMaxAge, // 旧文件最多保留90天
		Compress:  false,
		LocalTime: true,
	}

	return zapcore.AddSync(lumberJackLogger)
}
