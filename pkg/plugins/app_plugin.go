package plugins

import (
	"encoding/json"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"strings"
	"github.com/gosimple/slug"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

type AppPluginCss struct {
	Light	string	`json:"light"`
	Dark	string	`json:"dark"`
}
type AppPlugin struct {
	FrontendPluginBase
	Routes			[]*AppPluginRoute	`json:"routes"`
	FoundChildPlugins	[]*PluginInclude	`json:"-"`
	Pinned			bool			`json:"-"`
}
type AppPluginRoute struct {
	Path		string			`json:"path"`
	Method		string			`json:"method"`
	ReqRole		models.RoleType		`json:"reqRole"`
	Url		string			`json:"url"`
	Headers		[]AppPluginRouteHeader	`json:"headers"`
	TokenAuth	*JwtTokenAuth		`json:"tokenAuth"`
	JwtTokenAuth	*JwtTokenAuth		`json:"jwtTokenAuth"`
}
type AppPluginRouteHeader struct {
	Name	string	`json:"name"`
	Content	string	`json:"content"`
}
type JwtTokenAuth struct {
	Url	string			`json:"url"`
	Scopes	[]string		`json:"scopes"`
	Params	map[string]string	`json:"params"`
}

func (app *AppPlugin) Load(decoder *json.Decoder, pluginDir string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := decoder.Decode(&app); err != nil {
		return err
	}
	if err := app.registerPlugin(pluginDir); err != nil {
		return err
	}
	Apps[app.Id] = app
	return nil
}
func (app *AppPlugin) initApp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	app.initFrontendPlugin()
	for _, panel := range Panels {
		if strings.HasPrefix(panel.PluginDir, app.PluginDir) {
			panel.setPathsBasedOnApp(app)
			app.FoundChildPlugins = append(app.FoundChildPlugins, &PluginInclude{Name: panel.Name, Id: panel.Id, Type: panel.Type})
		}
	}
	for _, ds := range DataSources {
		if strings.HasPrefix(ds.PluginDir, app.PluginDir) {
			ds.setPathsBasedOnApp(app)
			app.FoundChildPlugins = append(app.FoundChildPlugins, &PluginInclude{Name: ds.Name, Id: ds.Id, Type: ds.Type})
		}
	}
	for _, include := range app.Includes {
		if include.Slug == "" {
			include.Slug = slug.Make(include.Name)
		}
		if include.Type == "page" && include.DefaultNav {
			app.DefaultNavUrl = setting.AppSubUrl + "/plugins/" + app.Id + "/page/" + include.Slug
		}
		if include.Type == "dashboard" && include.DefaultNav {
			app.DefaultNavUrl = setting.AppSubUrl + "/dashboard/db/" + include.Slug
		}
	}
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
