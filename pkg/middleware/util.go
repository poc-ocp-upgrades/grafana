package middleware

import (
	"strings"
	"github.com/go-macaron/gzip"
	"gopkg.in/macaron.v1"
)

func Gziper() macaron.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	macaronGziper := gzip.Gziper()
	return func(ctx *macaron.Context) {
		requestPath := ctx.Req.URL.RequestURI()
		if strings.HasPrefix(requestPath, "/api/datasources/proxy") {
			return
		}
		if strings.HasPrefix(requestPath, "/api/plugin-proxy/") {
			return
		}
		if strings.HasPrefix(requestPath, "/metrics") {
			return
		}
		ctx.Invoke(macaronGziper)
	}
}
