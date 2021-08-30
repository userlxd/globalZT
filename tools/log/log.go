package log

import (
	"globalZT/tools/config"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func init() {

	if config.Config.Mode == "dev" {
		log, err := zap.NewDevelopment()
		if err != nil {
			println("Logger Init Faild! Exit...")
			os.Exit(1)
		}
		Log = log.Sugar()
	} else {

		encoderConfig := zap.NewProductionEncoderConfig()
		// 修改时间编码器
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		// 在日志文件中使用大写字母记录日志级别
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		lumberJackLogger := &lumberjack.Logger{
			Filename:   config.Config.LogFile,
			MaxSize:    100,
			MaxBackups: 30,
			MaxAge:     30,
			Compress:   false,
		}

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(lumberJackLogger),
			zapcore.InfoLevel,
		)
		Log = zap.New(core, zap.AddCaller()).Sugar()
	}

}
