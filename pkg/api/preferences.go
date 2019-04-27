package api

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
)

func SetHomeDashboard(c *m.ReqContext, cmd m.SavePreferencesCommand) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cmd.UserId = c.UserId
	cmd.OrgId = c.OrgId
	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Failed to set home dashboard", err)
	}
	return Success("Home dashboard set")
}
func GetUserPreferences(c *m.ReqContext) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return getPreferencesFor(c.OrgId, c.UserId, 0)
}
func getPreferencesFor(orgID, userID, teamID int64) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	prefsQuery := m.GetPreferencesQuery{UserId: userID, OrgId: orgID, TeamId: teamID}
	if err := bus.Dispatch(&prefsQuery); err != nil {
		return Error(500, "Failed to get preferences", err)
	}
	dto := dtos.Prefs{Theme: prefsQuery.Result.Theme, HomeDashboardID: prefsQuery.Result.HomeDashboardId, Timezone: prefsQuery.Result.Timezone}
	return JSON(200, &dto)
}
func UpdateUserPreferences(c *m.ReqContext, dtoCmd dtos.UpdatePrefsCmd) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return updatePreferencesFor(c.OrgId, c.UserId, 0, &dtoCmd)
}
func updatePreferencesFor(orgID, userID, teamId int64, dtoCmd *dtos.UpdatePrefsCmd) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	saveCmd := m.SavePreferencesCommand{UserId: userID, OrgId: orgID, TeamId: teamId, Theme: dtoCmd.Theme, Timezone: dtoCmd.Timezone, HomeDashboardId: dtoCmd.HomeDashboardID}
	if err := bus.Dispatch(&saveCmd); err != nil {
		return Error(500, "Failed to save preferences", err)
	}
	return Success("Preferences updated")
}
func GetOrgPreferences(c *m.ReqContext) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return getPreferencesFor(c.OrgId, 0, 0)
}
func UpdateOrgPreferences(c *m.ReqContext, dtoCmd dtos.UpdatePrefsCmd) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return updatePreferencesFor(c.OrgId, 0, 0, &dtoCmd)
}
