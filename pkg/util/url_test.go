package util

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
)

func TestUrl(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	Convey("When joining two urls where right hand side is empty", t, func() {
		result := JoinUrlFragments("http://localhost:8080", "")
		So(result, ShouldEqual, "http://localhost:8080")
	})
	Convey("When joining two urls where right hand side is empty and lefthand side has a trailing slash", t, func() {
		result := JoinUrlFragments("http://localhost:8080/", "")
		So(result, ShouldEqual, "http://localhost:8080/")
	})
	Convey("When joining two urls where neither has a trailing slash", t, func() {
		result := JoinUrlFragments("http://localhost:8080", "api")
		So(result, ShouldEqual, "http://localhost:8080/api")
	})
	Convey("When joining two urls where lefthand side has a trailing slash", t, func() {
		result := JoinUrlFragments("http://localhost:8080/", "api")
		So(result, ShouldEqual, "http://localhost:8080/api")
	})
	Convey("When joining two urls where righthand side has preceding slash", t, func() {
		result := JoinUrlFragments("http://localhost:8080", "/api")
		So(result, ShouldEqual, "http://localhost:8080/api")
	})
	Convey("When joining two urls where righthand side has trailing slash", t, func() {
		result := JoinUrlFragments("http://localhost:8080", "api/")
		So(result, ShouldEqual, "http://localhost:8080/api/")
	})
	Convey("When joining two urls where lefthand side has a trailing slash and righthand side has preceding slash", t, func() {
		result := JoinUrlFragments("http://localhost:8080/", "/api/")
		So(result, ShouldEqual, "http://localhost:8080/api/")
	})
}
func TestNewUrlQueryReader(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	u, _ := url.Parse("http://www.abc.com/foo?bar=baz&bar2=baz2")
	uqr, _ := NewUrlQueryReader(u)
	Convey("when trying to retrieve the first query value", t, func() {
		result := uqr.Get("bar", "foodef")
		So(result, ShouldEqual, "baz")
	})
	Convey("when trying to retrieve the second query value", t, func() {
		result := uqr.Get("bar2", "foodef")
		So(result, ShouldEqual, "baz2")
	})
	Convey("when trying to retrieve from a non-existent key, the default value is returned", t, func() {
		result := uqr.Get("bar3", "foodef")
		So(result, ShouldEqual, "foodef")
	})
}
