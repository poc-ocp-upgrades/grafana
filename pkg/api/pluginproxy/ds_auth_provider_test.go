package pluginproxy

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDsAuthProvider(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	Convey("When interpolating string", t, func() {
		data := templateData{SecureJsonData: map[string]string{"Test": "0asd+asd"}}
		interpolated, err := interpolateString("{{.SecureJsonData.Test}}", data)
		So(err, ShouldBeNil)
		So(interpolated, ShouldEqual, "0asd+asd")
	})
}
