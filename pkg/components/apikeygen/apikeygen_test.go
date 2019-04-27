package apikeygen

import (
	"testing"
	"github.com/grafana/grafana/pkg/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApiKeyGen(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	Convey("When generating new api key", t, func() {
		result := New(12, "Cool key")
		So(result.ClientSecret, ShouldNotBeEmpty)
		So(result.HashedKey, ShouldNotBeEmpty)
		Convey("can decode key", func() {
			keyInfo, err := Decode(result.ClientSecret)
			So(err, ShouldBeNil)
			keyHashed := util.EncodePassword(keyInfo.Key, keyInfo.Name)
			So(keyHashed, ShouldEqual, result.HashedKey)
		})
	})
}
