package logger

import (
	"my-mall/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultLogger *zap.SugaredLogger

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	fileWriteSyncer := getFileLogWriter()

	var cores []zapcore.Core
	switch config.App.Env {
	case "test", "prod":
		// 测试环境和生产环境的日志输出到文件中
		cores = []zapcore.Core{zapcore.NewCore(encoder, fileWriteSyncer, zapcore.InfoLevel)}
	case "dev":
		// 开发环境同时向控制台和文件输出日志， Debug级别的日志也会被输出
		cores = []zapcore.Core{
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
		}
	}
	defaultLogger = zap.New(zapcore.NewTee(cores...)).Sugar()
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
