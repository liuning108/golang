package util

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

var log zerolog.Logger //创建日志变量

type Level int8 //定义日志等级

const (
	//调试模式
	DebugLevel Level = iota
	//正常输出
	InfoLevel
	//警告信息
	WarnLevel
	//错误信息
	ErrorLevel
	//严重错误信息
	FatalLevel
	//程序异常
	PanicLevel
	//没有等级
	NoLevel
	//禁用
	Disabled
)

//初始化
func init() {
	//初始设置为调试信息
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	//输出集
	var writers []io.Writer
	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	writers = append(writers, newRollingFile())
	//创建控制台输出对象,指定时间格式
	output := io.MultiWriter(writers...)

	//实例化日志对象
	log = zerolog.New(output).With().Timestamp().Logger()

	//设置调试等级
	SetLevel(DebugLevel)
}

func newRollingFile() io.Writer {
	return &lumberjack.Logger{
		Filename:   "log/log.log",
		MaxSize:    50, // megabytes
		MaxBackups: 5,
		MaxAge:     1,     //days
		Compress:   false, // disabled by default
	}
}

func SetLevel(l Level) {
	zerolog.SetGlobalLevel(zerolog.Level(l))
}

//输出调试日志信息
func Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

//正常输出日志信息
func Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

//输出警告日志信息
func Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

//输出错误日志信息
func Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

//输出异常日志信息
func Panicf(format string, v ...interface{}) {
	log.Panic().Msgf(format, v...)
}
