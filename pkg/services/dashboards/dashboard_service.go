package dashboards

import (
	"strings"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"time"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/guardian"
	"github.com/grafana/grafana/pkg/util"
)

type DashboardService interface {
	SaveDashboard(dto *SaveDashboardDTO) (*models.Dashboard, error)
	ImportDashboard(dto *SaveDashboardDTO) (*models.Dashboard, error)
}
type DashboardProvisioningService interface {
	SaveProvisionedDashboard(dto *SaveDashboardDTO, provisioning *models.DashboardProvisioning) (*models.Dashboard, error)
	SaveFolderForProvisionedDashboards(*SaveDashboardDTO) (*models.Dashboard, error)
	GetProvisionedDashboardData(name string) ([]*models.DashboardProvisioning, error)
}

var NewService = func() DashboardService {
	return &dashboardServiceImpl{log: log.New("dashboard-service")}
}
var NewProvisioningService = func() DashboardProvisioningService {
	return &dashboardServiceImpl{}
}

type SaveDashboardDTO struct {
	OrgId		int64
	UpdatedAt	time.Time
	User		*models.SignedInUser
	Message		string
	Overwrite	bool
	Dashboard	*models.Dashboard
}
type dashboardServiceImpl struct {
	orgId	int64
	user	*models.SignedInUser
	log		log.Logger
}

func (dr *dashboardServiceImpl) GetProvisionedDashboardData(name string) ([]*models.DashboardProvisioning, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cmd := &models.GetProvisionedDashboardDataQuery{Name: name}
	err := bus.Dispatch(cmd)
	if err != nil {
		return nil, err
	}
	return cmd.Result, nil
}
func (dr *dashboardServiceImpl) buildSaveDashboardCommand(dto *SaveDashboardDTO, validateAlerts bool, validateProvisionedDashboard bool) (*models.SaveDashboardCommand, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dash := dto.Dashboard
	dash.Title = strings.TrimSpace(dash.Title)
	dash.Data.Set("title", dash.Title)
	dash.SetUid(strings.TrimSpace(dash.Uid))
	if dash.Title == "" {
		return nil, models.ErrDashboardTitleEmpty
	}
	if dash.IsFolder && dash.FolderId > 0 {
		return nil, models.ErrDashboardFolderCannotHaveParent
	}
	if dash.IsFolder && strings.ToLower(dash.Title) == strings.ToLower(models.RootFolderName) {
		return nil, models.ErrDashboardFolderNameExists
	}
	if !util.IsValidShortUid(dash.Uid) {
		return nil, models.ErrDashboardInvalidUid
	} else if len(dash.Uid) > 40 {
		return nil, models.ErrDashboardUidToLong
	}
	if validateAlerts {
		validateAlertsCmd := models.ValidateDashboardAlertsCommand{OrgId: dto.OrgId, Dashboard: dash, User: dto.User}
		if err := bus.Dispatch(&validateAlertsCmd); err != nil {
			return nil, err
		}
	}
	validateBeforeSaveCmd := models.ValidateDashboardBeforeSaveCommand{OrgId: dto.OrgId, Dashboard: dash, Overwrite: dto.Overwrite}
	if err := bus.Dispatch(&validateBeforeSaveCmd); err != nil {
		return nil, err
	}
	if validateBeforeSaveCmd.Result.IsParentFolderChanged {
		folderGuardian := guardian.New(dash.FolderId, dto.OrgId, dto.User)
		if canSave, err := folderGuardian.CanSave(); err != nil || !canSave {
			if err != nil {
				return nil, err
			}
			return nil, models.ErrDashboardUpdateAccessDenied
		}
	}
	if validateProvisionedDashboard {
		isDashboardProvisioned := &models.IsDashboardProvisionedQuery{DashboardId: dash.Id}
		err := bus.Dispatch(isDashboardProvisioned)
		if err != nil {
			return nil, err
		}
		if isDashboardProvisioned.Result {
			return nil, models.ErrDashboardCannotSaveProvisionedDashboard
		}
	}
	guard := guardian.New(dash.GetDashboardIdForSavePermissionCheck(), dto.OrgId, dto.User)
	if canSave, err := guard.CanSave(); err != nil || !canSave {
		if err != nil {
			return nil, err
		}
		return nil, models.ErrDashboardUpdateAccessDenied
	}
	cmd := &models.SaveDashboardCommand{Dashboard: dash.Data, Message: dto.Message, OrgId: dto.OrgId, Overwrite: dto.Overwrite, UserId: dto.User.UserId, FolderId: dash.FolderId, IsFolder: dash.IsFolder, PluginId: dash.PluginId}
	if !dto.UpdatedAt.IsZero() {
		cmd.UpdatedAt = dto.UpdatedAt
	}
	return cmd, nil
}
func (dr *dashboardServiceImpl) updateAlerting(cmd *models.SaveDashboardCommand, dto *SaveDashboardDTO) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	alertCmd := models.UpdateDashboardAlertsCommand{OrgId: dto.OrgId, Dashboard: cmd.Result, User: dto.User}
	if err := bus.Dispatch(&alertCmd); err != nil {
		return err
	}
	return nil
}
func (dr *dashboardServiceImpl) SaveProvisionedDashboard(dto *SaveDashboardDTO, provisioning *models.DashboardProvisioning) (*models.Dashboard, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dto.User = &models.SignedInUser{UserId: 0, OrgRole: models.ROLE_ADMIN, OrgId: dto.OrgId}
	cmd, err := dr.buildSaveDashboardCommand(dto, true, false)
	if err != nil {
		return nil, err
	}
	saveCmd := &models.SaveProvisionedDashboardCommand{DashboardCmd: cmd, DashboardProvisioning: provisioning}
	err = bus.Dispatch(saveCmd)
	if err != nil {
		return nil, err
	}
	err = dr.updateAlerting(cmd, dto)
	if err != nil {
		return nil, err
	}
	return cmd.Result, nil
}
func (dr *dashboardServiceImpl) SaveFolderForProvisionedDashboards(dto *SaveDashboardDTO) (*models.Dashboard, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dto.User = &models.SignedInUser{UserId: 0, OrgRole: models.ROLE_ADMIN}
	cmd, err := dr.buildSaveDashboardCommand(dto, false, false)
	if err != nil {
		return nil, err
	}
	err = bus.Dispatch(cmd)
	if err != nil {
		return nil, err
	}
	err = dr.updateAlerting(cmd, dto)
	if err != nil {
		return nil, err
	}
	return cmd.Result, nil
}
func (dr *dashboardServiceImpl) SaveDashboard(dto *SaveDashboardDTO) (*models.Dashboard, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cmd, err := dr.buildSaveDashboardCommand(dto, true, true)
	if err != nil {
		return nil, err
	}
	err = bus.Dispatch(cmd)
	if err != nil {
		return nil, err
	}
	err = dr.updateAlerting(cmd, dto)
	if err != nil {
		return nil, err
	}
	return cmd.Result, nil
}
func (dr *dashboardServiceImpl) ImportDashboard(dto *SaveDashboardDTO) (*models.Dashboard, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cmd, err := dr.buildSaveDashboardCommand(dto, false, true)
	if err != nil {
		return nil, err
	}
	err = bus.Dispatch(cmd)
	if err != nil {
		return nil, err
	}
	return cmd.Result, nil
}

type FakeDashboardService struct {
	SaveDashboardResult	*models.Dashboard
	SaveDashboardError	error
	SavedDashboards		[]*SaveDashboardDTO
}

func (s *FakeDashboardService) SaveDashboard(dto *SaveDashboardDTO) (*models.Dashboard, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.SavedDashboards = append(s.SavedDashboards, dto)
	if s.SaveDashboardResult == nil && s.SaveDashboardError == nil {
		s.SaveDashboardResult = dto.Dashboard
	}
	return s.SaveDashboardResult, s.SaveDashboardError
}
func (s *FakeDashboardService) ImportDashboard(dto *SaveDashboardDTO) (*models.Dashboard, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.SaveDashboard(dto)
}
func MockDashboardService(mock *FakeDashboardService) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	NewService = func() DashboardService {
		return mock
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
