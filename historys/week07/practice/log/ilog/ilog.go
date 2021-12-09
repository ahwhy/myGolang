package ilog

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger(level string) {
	// 设置日志等级
	switch level {
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	}
	
	// 设置日志输出路径
	Logger.SetOutput(os.Stdout)
	// 设置filename
	Logger.SetReportCaller(true)
	// 设置格式化文本
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
	})
}
