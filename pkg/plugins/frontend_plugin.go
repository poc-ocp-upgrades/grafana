package plugins

import (
	"net/url"
	"path"
	"strings"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

type FrontendPluginBase struct{ PluginBase }

func (fp *FrontendPluginBase) initFrontendPlugin() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if isExternalPlugin(fp.PluginDir) {
		StaticRoutes = append(StaticRoutes, &PluginStaticRoute{Directory: fp.PluginDir, PluginId: fp.Id})
	}
	fp.handleModuleDefaults()
	fp.Info.Logos.Small = getPluginLogoUrl(fp.Type, fp.Info.Logos.Small, fp.BaseUrl)
	fp.Info.Logos.Large = getPluginLogoUrl(fp.Type, fp.Info.Logos.Large, fp.BaseUrl)
	for i := 0; i < len(fp.Info.Screenshots); i++ {
		fp.Info.Screenshots[i].Path = evalRelativePluginUrlPath(fp.Info.Screenshots[i].Path, fp.BaseUrl)
	}
}
func getPluginLogoUrl(pluginType, path, baseUrl string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if path == "" {
		return "public/img/icn-" + pluginType + ".svg"
	}
	return evalRelativePluginUrlPath(path, baseUrl)
}
func (fp *FrontendPluginBase) setPathsBasedOnApp(app *AppPlugin) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	appSubPath := strings.Replace(strings.Replace(fp.PluginDir, app.PluginDir, "", 1), "\\", "/", -1)
	fp.IncludedInAppId = app.Id
	fp.BaseUrl = app.BaseUrl
	if isExternalPlugin(app.PluginDir) {
		fp.Module = util.JoinUrlFragments("plugins/"+app.Id, appSubPath) + "/module"
	} else {
		fp.Module = util.JoinUrlFragments("app/plugins/app/"+app.Id, appSubPath) + "/module"
	}
}
func (fp *FrontendPluginBase) handleModuleDefaults() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if isExternalPlugin(fp.PluginDir) {
		fp.Module = path.Join("plugins", fp.Id, "module")
		fp.BaseUrl = path.Join("public/plugins", fp.Id)
		return
	}
	fp.IsCorePlugin = true
	fp.Module = path.Join("app/plugins", fp.Type, fp.Id, "module")
	fp.BaseUrl = path.Join("public/app/plugins", fp.Type, fp.Id)
}
func isExternalPlugin(pluginDir string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return !strings.Contains(pluginDir, setting.StaticRootPath)
}
func evalRelativePluginUrlPath(pathStr string, baseUrl string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if pathStr == "" {
		return ""
	}
	u, _ := url.Parse(pathStr)
	if u.IsAbs() {
		return pathStr
	}
	return path.Join(baseUrl, pathStr)
}
