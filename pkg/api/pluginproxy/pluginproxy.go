package pluginproxy

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/util"
)

type templateData struct {
	JsonData	map[string]interface{}
	SecureJsonData	map[string]string
}

func getHeaders(route *plugins.AppPluginRoute, orgId int64, appID string) (http.Header, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := http.Header{}
	query := m.GetPluginSettingByIdQuery{OrgId: orgId, PluginId: appID}
	if err := bus.Dispatch(&query); err != nil {
		return nil, err
	}
	data := templateData{JsonData: query.Result.JsonData, SecureJsonData: query.Result.SecureJsonData.Decrypt()}
	err := addHeaders(&result, route, data)
	return result, err
}
func NewApiPluginProxy(ctx *m.ReqContext, proxyPath string, route *plugins.AppPluginRoute, appID string) *httputil.ReverseProxy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	targetURL, _ := url.Parse(route.Url)
	director := func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		req.URL.Path = util.JoinUrlFragments(targetURL.Path, proxyPath)
		req.Header.Del("Cookie")
		req.Header.Del("Set-Cookie")
		req.Header.Del("X-Forwarded-Host")
		req.Header.Del("X-Forwarded-Port")
		req.Header.Del("X-Forwarded-Proto")
		if req.RemoteAddr != "" {
			remoteAddr, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				remoteAddr = req.RemoteAddr
			}
			if req.Header.Get("X-Forwarded-For") != "" {
				req.Header.Set("X-Forwarded-For", req.Header.Get("X-Forwarded-For")+", "+remoteAddr)
			} else {
				req.Header.Set("X-Forwarded-For", remoteAddr)
			}
		}
		ctxJson, err := json.Marshal(ctx.SignedInUser)
		if err != nil {
			ctx.JsonApiErr(500, "failed to marshal context to json.", err)
			return
		}
		req.Header.Add("X-Grafana-Context", string(ctxJson))
		if len(route.Headers) > 0 {
			headers, err := getHeaders(route, ctx.OrgId, appID)
			if err != nil {
				ctx.JsonApiErr(500, "Could not generate plugin route header", err)
				return
			}
			for key, value := range headers {
				log.Trace("setting key %v value <redacted>", key)
				req.Header.Set(key, value[0])
			}
		}
	}
	return &httputil.ReverseProxy{Director: director}
}
