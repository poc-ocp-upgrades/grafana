package api

import (
	"github.com/grafana/grafana/pkg/api/pluginproxy"
	"github.com/grafana/grafana/pkg/metrics"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
)

func (hs *HTTPServer) ProxyDataSourceRequest(c *m.ReqContext) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.TimeRequest(metrics.M_DataSource_ProxyReq_Timer)
	dsId := c.ParamsInt64(":id")
	ds, err := hs.DatasourceCache.GetDatasource(dsId, c.SignedInUser, c.SkipCache)
	if err != nil {
		if err == m.ErrDataSourceAccessDenied {
			c.JsonApiErr(403, "Access denied to datasource", err)
			return
		}
		c.JsonApiErr(500, "Unable to load datasource meta data", err)
		return
	}
	plugin, ok := plugins.DataSources[ds.Type]
	if !ok {
		c.JsonApiErr(500, "Unable to find datasource plugin", err)
		return
	}
	proxyPath := ensureProxyPathTrailingSlash(c.Req.URL.Path, c.Params("*"))
	proxy := pluginproxy.NewDataSourceProxy(ds, plugin, c, proxyPath)
	proxy.HandleRequest()
}
func ensureProxyPathTrailingSlash(originalPath, proxyPath string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(proxyPath) > 1 {
		if originalPath[len(originalPath)-1] == '/' && proxyPath[len(proxyPath)-1] != '/' {
			return proxyPath + "/"
		}
	}
	return proxyPath
}
