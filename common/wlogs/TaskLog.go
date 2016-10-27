package wlogs

import (
	"fmt"
	"log"
)

func Info(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Infofln(format string, args ...interface{}) {
	fmt.Sprintf(format, args...)
}

func Infof(w io.Writer, args ...interface{}) {
	fmt.Fprinln(w, args...)
}

func Error(format string, args ...interface{}) {
	fmt.Errorf(format, args...)
}
