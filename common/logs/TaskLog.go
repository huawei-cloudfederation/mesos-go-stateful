package logs

import (
	"flag"
	"github.com/golang/glog"
)

//Printf logs to the console
func Printf(format string, args ...interface{}) {
	flag.Parse()
	glog.Infof(format, args...)
	flag.Lookup("logtostderr").Value.Set("true")
}

//Error logs the error to the console
func Error(format string, args ...interface{}) {
	flag.Parse()
	glog.Errorf(format, args...)
	flag.Lookup("logtostderr").Value.Set("true")
}

//Fatal logs the fatal error to the console
func Fatal(args ...interface{}) {
	flag.Parse()
	glog.Fatal(args...)
	flag.Lookup("logtostderr").Value.Set("true")
}

//FatalInfo logs the fatal error to the console
func FatalInfo(format string, args ...interface{}) {
	flag.Parse()
	glog.Fatalf(format, args...)
	flag.Lookup("logtostderr").Value.Set("true")
}

//Println logs to the console
func Println(args ...interface{}) {
	flag.Parse()
	glog.Infoln(args...)
	flag.Lookup("logtostderr").Value.Set("true")
}
