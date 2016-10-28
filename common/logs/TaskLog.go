package logs

import (
	"glog"
)

func Printf(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

func Fatalf(args ...interface{}) {
	glog.Fatal(args...)
}

func Println(args ...interface{}) {
	glog.Infoln(args...)
}


func Error(format string, args ...interface{}) {
	glog.Errorf(format, args...)
}
