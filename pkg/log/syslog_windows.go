package log

import (
	"github.com/inconshreveable/log15"
	"gopkg.in/ini.v1"
)

type SysLogHandler struct{}

func NewSyslog(sec *ini.Section, format log15.Format) *SysLogHandler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &SysLogHandler{}
}
func (sw *SysLogHandler) Log(r *log15.Record) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (sw *SysLogHandler) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
