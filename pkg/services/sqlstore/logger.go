package sqlstore

import (
	"fmt"
	glog "github.com/grafana/grafana/pkg/log"
	"github.com/go-xorm/core"
)

type XormLogger struct {
	grafanaLog	glog.Logger
	level		glog.Lvl
	showSQL		bool
}

func NewXormLogger(level glog.Lvl, grafanaLog glog.Logger) *XormLogger {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &XormLogger{grafanaLog: grafanaLog, level: level, showSQL: true}
}
func (s *XormLogger) Error(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.level <= glog.LvlError {
		s.grafanaLog.Error(fmt.Sprint(v...))
	}
}
func (s *XormLogger) Errorf(format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.level <= glog.LvlError {
		s.grafanaLog.Error(fmt.Sprintf(format, v...))
	}
}
func (s *XormLogger) Debug(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.level <= glog.LvlDebug {
		s.grafanaLog.Debug(fmt.Sprint(v...))
	}
}
func (s *XormLogger) Debugf(format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.level <= glog.LvlDebug {
		s.grafanaLog.Debug(fmt.Sprintf(format, v...))
	}
}
func (s *XormLogger) Info(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.level <= glog.LvlInfo {
		s.grafanaLog.Info(fmt.Sprint(v...))
	}
}
func (s *XormLogger) Infof(format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.level <= glog.LvlInfo {
		s.grafanaLog.Info(fmt.Sprintf(format, v...))
	}
}
func (s *XormLogger) Warn(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.level <= glog.LvlWarn {
		s.grafanaLog.Warn(fmt.Sprint(v...))
	}
}
func (s *XormLogger) Warnf(format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.level <= glog.LvlWarn {
		s.grafanaLog.Warn(fmt.Sprintf(format, v...))
	}
}
func (s *XormLogger) Level() core.LogLevel {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch s.level {
	case glog.LvlError:
		return core.LOG_ERR
	case glog.LvlWarn:
		return core.LOG_WARNING
	case glog.LvlInfo:
		return core.LOG_INFO
	case glog.LvlDebug:
		return core.LOG_DEBUG
	default:
		return core.LOG_ERR
	}
}
func (s *XormLogger) SetLevel(l core.LogLevel) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (s *XormLogger) ShowSQL(show ...bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.grafanaLog.Error("ShowSQL", "show", "show")
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}
func (s *XormLogger) IsShowSQL() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.showSQL
}
