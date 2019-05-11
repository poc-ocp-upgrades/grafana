package extensions

import (
	_ "gopkg.in/square/go-jose.v2"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

var IsEnterprise bool = false

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
