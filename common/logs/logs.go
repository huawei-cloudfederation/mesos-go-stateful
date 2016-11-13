package logs

import (
	"fmt"
	"github.com/golang/glog"
)

//Printf logs to the console
func Printf(format string, args ...interface{}) {
	glog.InfoDepth(1, fmt.Sprintf(format, args...))
}

//Error logs the error to the console
func Error(format string, args ...interface{}) {
	glog.ErrorDepth(1, fmt.Sprintf(format, args...))
}

//Fatal logs the fatal error to the console
func Fatal(args ...interface{}) {
	glog.FatalDepth(1, args...)
}

//FatalInfo logs the fatal error to the console
func FatalInfo(format string, args ...interface{}) {
	glog.FatalDepth(1, fmt.Sprintf(format, args...))
}

//Println logs to the console
func Println(args ...interface{}) {
	glog.InfoDepth(1, args...)
}
