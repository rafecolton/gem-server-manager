package log

import (
	stdlog "log"
	"os"
)

type Log interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

type GsmLogger struct {
	Log
}

func (l *GsmLogger) IsBoolFlag() bool {
	return true
}

func (l *GsmLogger) Initialize() {
	if l.Log == nil {
		l.Log = &nullLogger{}
	}
}

//to satisfy the flag Value interface
func (l *GsmLogger) String() string {
	return ""
}

//to satisfy the flag Value interface
func (l *GsmLogger) Set(value string) error {
	if value == "true" {
		l.Log = stdlog.New(os.Stderr, "[gsm] ", stdlog.LstdFlags)
	} else {
		l.Log = &nullLogger{}
	}
	return nil
}
