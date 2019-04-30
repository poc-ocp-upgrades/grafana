package plugins

import (
	"log"
	glog "github.com/grafana/grafana/pkg/log"
	hclog "github.com/hashicorp/go-hclog"
)

type LogWrapper struct{ Logger glog.Logger }

func (lw LogWrapper) Trace(msg string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lw.Logger.Debug(msg, args...)
}
func (lw LogWrapper) Debug(msg string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lw.Logger.Debug(msg, args...)
}
func (lw LogWrapper) Info(msg string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lw.Logger.Info(msg, args...)
}
func (lw LogWrapper) Warn(msg string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lw.Logger.Warn(msg, args...)
}
func (lw LogWrapper) Error(msg string, args ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lw.Logger.Error(msg, args...)
}
func (lw LogWrapper) IsTrace() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (lw LogWrapper) IsDebug() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (lw LogWrapper) IsInfo() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (lw LogWrapper) IsWarn() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (lw LogWrapper) IsError() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (lw LogWrapper) With(args ...interface{}) hclog.Logger {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return LogWrapper{Logger: lw.Logger.New(args...)}
}
func (lw LogWrapper) Named(name string) hclog.Logger {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return LogWrapper{Logger: lw.Logger.New()}
}
func (lw LogWrapper) ResetNamed(name string) hclog.Logger {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return LogWrapper{Logger: lw.Logger.New()}
}
func (lw LogWrapper) StandardLogger(ops *hclog.StandardLoggerOptions) *log.Logger {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
