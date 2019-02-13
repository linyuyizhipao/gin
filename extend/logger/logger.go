package logger

import (
	"os"
	"strings"

	"test/extend/conf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup 日志初始化设置
func Setup() {
	level := strings.ToLower(conf.LoggerConf.Level) //获取配置文件中配置的日志异常级别，并转全小写
	//根据不同的日志级别设置zerolog组件 打印的日志格式
	switch  level{
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	//是否在终端输出日志，及输出格式
	if conf.LoggerConf.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: !conf.LoggerConf.Color,
		})
	}
}
