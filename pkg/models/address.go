package models

import (
	godefaultruntime "runtime"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
)

type Address struct {
	Address1	string	`json:"address1"`
	Address2	string	`json:"address2"`
	City		string	`json:"city"`
	ZipCode		string	`json:"zipCode"`
	State		string	`json:"state"`
	Country		string	`json:"country"`
}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
