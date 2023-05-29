package logger

import (
	"fmt"
	"log"
	"os"
)

const flag = log.Ldate | log.Ltime | log.Lshortfile

type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func New(prefix string) *Logger {
	self := Logger{
		info:  log.New(os.Stdout, fmt.Sprintf("%s %s ", prefix, "info"), flag),
		warn:  log.New(os.Stdout, fmt.Sprintf("%s %s ", prefix, "warn"), flag),
		error: log.New(os.Stderr, fmt.Sprintf("%s %s ", prefix, "error"), flag),
	}

	return &self
}

func (self *Logger) Info(v ...any) {
	self.info.Print(v...)
}

func (self *Logger) Infof(format string, v ...any) {
	self.info.Printf(format, v...)
}

func (self *Logger) Infoln(v ...any) {
	self.info.Println(v...)
}

func (self *Logger) Warn(v ...any) {
	self.warn.Print(v...)
}

func (self *Logger) Warnf(format string, v ...any) {
	self.warn.Printf(format, v...)
}

func (self *Logger) Warnln(v ...any) {
	self.warn.Println(v...)
}

func (self *Logger) Error(v ...any) {
	self.error.Fatal(v...)
}

func (self *Logger) Errorf(format string, v ...any) {
	self.error.Fatalf(format, v...)
}

func (self *Logger) Errorln(v ...any) {
	self.error.Fatalln(v...)
}
