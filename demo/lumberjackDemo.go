package main

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "log.log",
		MaxSize:    50, // megabytes
		MaxBackups: 5,
		MaxAge:     1,     //days
		Compress:   false, // disabled by default
	})

}
