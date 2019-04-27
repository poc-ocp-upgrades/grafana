package alerting

import (
	"errors"
	"fmt"
	"time"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
)

type DashAlertExtractor struct {
	User	*m.SignedInUser
	Dash	*m.Dashboard
	OrgID	int64
	log	log.Logger
}

func NewDashAlertExtractor(dash *m.Dashboard, orgID int64, user *m.SignedInUser) *DashAlertExtractor {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &DashAlertExtractor{User: user, Dash: dash, OrgID: orgID, log: log.New("alerting.extractor")}
}
func (e *DashAlertExtractor) lookupDatasourceID(dsName string) (*m.DataSource, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if dsName == "" {
		query := &m.GetDataSourcesQuery{OrgId: e.OrgID}
		if err := bus.Dispatch(query); err != nil {
			return nil, err
		}
		for _, ds := range query.Result {
			if ds.IsDefault {
				return ds, nil
			}
		}
	} else {
		query := &m.GetDataSourceByNameQuery{Name: dsName, OrgId: e.OrgID}
		if err := bus.Dispatch(query); err != nil {
			return nil, err
		}
		return query.Result, nil
	}
	return nil, errors.New("Could not find datasource id for " + dsName)
}
func findPanelQueryByRefID(panel *simplejson.Json, refID string) *simplejson.Json {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, targetsObj := range panel.Get("targets").MustArray() {
		target := simplejson.NewFromAny(targetsObj)
		if target.Get("refId").MustString() == refID {
			return target
		}
	}
	return nil
}
func copyJSON(in *simplejson.Json) (*simplejson.Json, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rawJSON, err := in.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(rawJSON)
}
func (e *DashAlertExtractor) getAlertFromPanels(jsonWithPanels *simplejson.Json, validateAlertFunc func(*m.Alert) bool) ([]*m.Alert, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	alerts := make([]*m.Alert, 0)
	for _, panelObj := range jsonWithPanels.Get("panels").MustArray() {
		panel := simplejson.NewFromAny(panelObj)
		collapsedJSON, collapsed := panel.CheckGet("collapsed")
		if collapsed && collapsedJSON.MustBool() {
			alertSlice, err := e.getAlertFromPanels(panel, validateAlertFunc)
			if err != nil {
				return nil, err
			}
			alerts = append(alerts, alertSlice...)
			continue
		}
		jsonAlert, hasAlert := panel.CheckGet("alert")
		if !hasAlert {
			continue
		}
		panelID, err := panel.Get("id").Int64()
		if err != nil {
			return nil, ValidationError{Reason: "A numeric panel id property is missing"}
		}
		enabled, hasEnabled := jsonAlert.CheckGet("enabled")
		if hasEnabled && !enabled.MustBool() {
			continue
		}
		frequency, err := getTimeDurationStringToSeconds(jsonAlert.Get("frequency").MustString())
		if err != nil {
			return nil, ValidationError{Reason: err.Error()}
		}
		rawFor := jsonAlert.Get("for").MustString()
		var forValue time.Duration
		if rawFor != "" {
			forValue, err = time.ParseDuration(rawFor)
			if err != nil {
				return nil, ValidationError{Reason: "Could not parse for"}
			}
		}
		alert := &m.Alert{DashboardId: e.Dash.Id, OrgId: e.OrgID, PanelId: panelID, Id: jsonAlert.Get("id").MustInt64(), Name: jsonAlert.Get("name").MustString(), Handler: jsonAlert.Get("handler").MustInt64(), Message: jsonAlert.Get("message").MustString(), Frequency: frequency, For: forValue}
		for _, condition := range jsonAlert.Get("conditions").MustArray() {
			jsonCondition := simplejson.NewFromAny(condition)
			jsonQuery := jsonCondition.Get("query")
			queryRefID := jsonQuery.Get("params").MustArray()[0].(string)
			panelQuery := findPanelQueryByRefID(panel, queryRefID)
			if panelQuery == nil {
				reason := fmt.Sprintf("Alert on PanelId: %v refers to query(%s) that cannot be found", alert.PanelId, queryRefID)
				return nil, ValidationError{Reason: reason}
			}
			dsName := ""
			if panelQuery.Get("datasource").MustString() != "" {
				dsName = panelQuery.Get("datasource").MustString()
			} else if panel.Get("datasource").MustString() != "" {
				dsName = panel.Get("datasource").MustString()
			}
			datasource, err := e.lookupDatasourceID(dsName)
			if err != nil {
				e.log.Debug("Error looking up datasource", "error", err)
				return nil, ValidationError{Reason: fmt.Sprintf("Data source used by alert rule not found, alertName=%v, datasource=%s", alert.Name, dsName)}
			}
			dsFilterQuery := m.DatasourcesPermissionFilterQuery{User: e.User, Datasources: []*m.DataSource{datasource}}
			if err := bus.Dispatch(&dsFilterQuery); err != nil {
				if err != bus.ErrHandlerNotFound {
					return nil, err
				}
			} else {
				if len(dsFilterQuery.Result) == 0 {
					return nil, m.ErrDataSourceAccessDenied
				}
			}
			jsonQuery.SetPath([]string{"datasourceId"}, datasource.Id)
			if interval, err := panel.Get("interval").String(); err == nil {
				panelQuery.Set("interval", interval)
			}
			jsonQuery.Set("model", panelQuery.Interface())
		}
		alert.Settings = jsonAlert
		_, err = NewRuleFromDBAlert(alert)
		if err != nil {
			return nil, err
		}
		if !validateAlertFunc(alert) {
			return nil, ValidationError{Reason: fmt.Sprintf("Panel id is not correct, alertName=%v, panelId=%v", alert.Name, alert.PanelId)}
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}
func validateAlertRule(alert *m.Alert) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return alert.ValidToSave()
}
func (e *DashAlertExtractor) GetAlerts() ([]*m.Alert, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.extractAlerts(validateAlertRule)
}
func (e *DashAlertExtractor) extractAlerts(validateFunc func(alert *m.Alert) bool) ([]*m.Alert, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dashboardJSON, err := copyJSON(e.Dash.Data)
	if err != nil {
		return nil, err
	}
	alerts := make([]*m.Alert, 0)
	rows := dashboardJSON.Get("rows").MustArray()
	if len(rows) > 0 {
		for _, rowObj := range rows {
			row := simplejson.NewFromAny(rowObj)
			a, err := e.getAlertFromPanels(row, validateFunc)
			if err != nil {
				return nil, err
			}
			alerts = append(alerts, a...)
		}
	} else {
		a, err := e.getAlertFromPanels(dashboardJSON, validateFunc)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, a...)
	}
	e.log.Debug("Extracted alerts from dashboard", "alertCount", len(alerts))
	return alerts, nil
}
func (e *DashAlertExtractor) ValidateAlerts() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := e.extractAlerts(func(alert *m.Alert) bool {
		return alert.OrgId != 0 && alert.PanelId != 0
	})
	return err
}
