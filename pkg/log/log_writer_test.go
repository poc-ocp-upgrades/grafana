package log

import (
	"testing"
	"github.com/inconshreveable/log15"
	. "github.com/smartystreets/goconvey/convey"
)

type FakeLogger struct {
	debug	string
	info	string
	warn	string
	err	string
	crit	string
}

func (f *FakeLogger) New(ctx ...interface{}) log15.Logger {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (f *FakeLogger) Debug(msg string, ctx ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.debug = msg
}
func (f *FakeLogger) Info(msg string, ctx ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.info = msg
}
func (f *FakeLogger) Warn(msg string, ctx ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.warn = msg
}
func (f *FakeLogger) Error(msg string, ctx ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.err = msg
}
func (f *FakeLogger) Crit(msg string, ctx ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.crit = msg
}
func (f *FakeLogger) GetHandler() log15.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (f *FakeLogger) SetHandler(l log15.Handler) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func TestLogWriter(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	Convey("When writing to a LogWriter", t, func() {
		Convey("Should write using the correct level [crit]", func() {
			fake := &FakeLogger{}
			crit := NewLogWriter(fake, LvlCrit, "")
			n, err := crit.Write([]byte("crit"))
			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(fake.crit, ShouldEqual, "crit")
		})
		Convey("Should write using the correct level [error]", func() {
			fake := &FakeLogger{}
			crit := NewLogWriter(fake, LvlError, "")
			n, err := crit.Write([]byte("error"))
			So(n, ShouldEqual, 5)
			So(err, ShouldBeNil)
			So(fake.err, ShouldEqual, "error")
		})
		Convey("Should write using the correct level [warn]", func() {
			fake := &FakeLogger{}
			crit := NewLogWriter(fake, LvlWarn, "")
			n, err := crit.Write([]byte("warn"))
			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(fake.warn, ShouldEqual, "warn")
		})
		Convey("Should write using the correct level [info]", func() {
			fake := &FakeLogger{}
			crit := NewLogWriter(fake, LvlInfo, "")
			n, err := crit.Write([]byte("info"))
			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(fake.info, ShouldEqual, "info")
		})
		Convey("Should write using the correct level [debug]", func() {
			fake := &FakeLogger{}
			crit := NewLogWriter(fake, LvlDebug, "")
			n, err := crit.Write([]byte("debug"))
			So(n, ShouldEqual, 5)
			So(err, ShouldBeNil)
			So(fake.debug, ShouldEqual, "debug")
		})
		Convey("Should prefix the output with the prefix", func() {
			fake := &FakeLogger{}
			crit := NewLogWriter(fake, LvlDebug, "prefix")
			n, err := crit.Write([]byte("debug"))
			So(n, ShouldEqual, 5)
			So(err, ShouldBeNil)
			So(fake.debug, ShouldEqual, "prefixdebug")
		})
	})
}
