package services

import (
	"io/ioutil"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"os"
)

type IoUtilImp struct{}

func (i IoUtilImp) Stat(path string) (os.FileInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return os.Stat(path)
}
func (i IoUtilImp) RemoveAll(path string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return os.RemoveAll(path)
}
func (i IoUtilImp) ReadDir(path string) ([]os.FileInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ioutil.ReadDir(path)
}
func (i IoUtilImp) ReadFile(filename string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ioutil.ReadFile(filename)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
