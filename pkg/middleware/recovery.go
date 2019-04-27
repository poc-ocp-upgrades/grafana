package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"gopkg.in/macaron.v1"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

var (
	dunno		= []byte("???")
	centerDot	= []byte("·")
	dot		= []byte(".")
	slash		= []byte("/")
)

func stack(skip int) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}
func source(lines [][]byte, n int) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}
func function(pc uintptr) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
func Recovery() macaron.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(c *macaron.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				panicLogger := log.Root
				if ctx, ok := c.Data["ctx"]; ok {
					ctxTyped := ctx.(*m.ReqContext)
					panicLogger = ctxTyped.Logger
				}
				panicLogger.Error("Request error", "error", err, "stack", string(stack))
				c.Data["Title"] = "Server Error"
				c.Data["AppSubUrl"] = setting.AppSubUrl
				c.Data["Theme"] = setting.DefaultTheme
				if setting.Env == setting.DEV {
					if theErr, ok := err.(error); ok {
						c.Data["Title"] = theErr.Error()
					}
					c.Data["ErrorMsg"] = string(stack)
				}
				ctx, ok := c.Data["ctx"].(*m.ReqContext)
				if ok && ctx.IsApiRequest() {
					resp := make(map[string]interface{})
					resp["message"] = "Internal Server Error - Check the Grafana server logs for the detailed error message."
					if c.Data["ErrorMsg"] != nil {
						resp["error"] = fmt.Sprintf("%v - %v", c.Data["Title"], c.Data["ErrorMsg"])
					} else {
						resp["error"] = c.Data["Title"]
					}
					c.JSON(500, resp)
				} else {
					c.HTML(500, setting.ERR_TEMPLATE_NAME)
				}
			}
		}()
		c.Next()
	}
}
