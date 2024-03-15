package app

import (
	"io"
	"log"
	"os"
)

type logger struct {
	debug *log.Logger
	info  *log.Logger
	err   *log.Logger
}

func newLogger(debugEnabled bool) *logger {
	dbgWriter := io.Discard
	if debugEnabled {
		dbgWriter = os.Stderr
	}

	return &logger{
		debug: log.New(dbgWriter, "[DBG] ", log.LstdFlags),
		info:  log.New(os.Stderr, "[INF] ", log.LstdFlags),
		err:   log.New(os.Stderr, "[ERR]", log.LstdFlags),
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}

func (l *logger) Debug(v ...interface{}) {
	l.debug.Print(v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.info.Print(v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.err.Printf(format, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.err.Print(v...)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.err.Fatalf(format, v...)
}

func (l *logger) Fatal(v ...interface{}) {
	l.err.Fatal(v...)
}

func (l *logger) Panicf(format string, v ...interface{}) {
	l.err.Panicf(format, v...)
}

func (l *logger) Panic(v ...interface{}) {
	l.err.Panic(v...)
}
