package middleware

import (
	"net/http"
	"time"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/macaron.v1"
)

func Logger() macaron.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		start := time.Now()
		c.Data["perfmon.start"] = start
		rw := res.(macaron.ResponseWriter)
		c.Next()
		timeTakenMs := time.Since(start) / time.Millisecond
		if timer, ok := c.Data["perfmon.timer"]; ok {
			timerTyped := timer.(prometheus.Summary)
			timerTyped.Observe(float64(timeTakenMs))
		}
		status := rw.Status()
		if status == 200 || status == 304 {
			if !setting.RouterLogging {
				return
			}
		}
		if ctx, ok := c.Data["ctx"]; ok {
			ctxTyped := ctx.(*m.ReqContext)
			if status == 500 {
				ctxTyped.Logger.Error("Request Completed", "method", req.Method, "path", req.URL.Path, "status", status, "remote_addr", c.RemoteAddr(), "time_ms", int64(timeTakenMs), "size", rw.Size(), "referer", req.Referer())
			} else {
				ctxTyped.Logger.Info("Request Completed", "method", req.Method, "path", req.URL.Path, "status", status, "remote_addr", c.RemoteAddr(), "time_ms", int64(timeTakenMs), "size", rw.Size(), "referer", req.Referer())
			}
		}
	}
}
