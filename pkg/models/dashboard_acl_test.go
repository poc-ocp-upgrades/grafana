package models

import (
	"testing"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDashboardAclModel(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	Convey("When printing a PermissionType", t, func() {
		view := PERMISSION_VIEW
		printed := fmt.Sprint(view)
		Convey("Should output a friendly name", func() {
			So(printed, ShouldEqual, "View")
		})
	})
}
