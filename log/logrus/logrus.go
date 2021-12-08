package logrus

import (
	"os"

	logger "github.com/sirupsen/logrus"
)

var Logger = logger.New()

func InitLogger(level string) {
	// 设置日志等级
	switch level {
	case "info":
		Logger.SetLevel(logger.InfoLevel)
	case "debug":
		Logger.SetLevel(logger.DebugLevel)
	case "warning":
		Logger.SetLevel(logger.WarnLevel)
	case "error":
		Logger.SetLevel(logger.ErrorLevel)
	}

	// 设置日志输出路径
	Logger.SetOutput(os.Stdout)
	// 设置filename
	Logger.SetReportCaller(true)
	// 设置格式化文本
	Logger.SetFormatter(&logger.TextFormatter{
		TimestampFormat: "2006/1/2 15:04:05",
	})
}
