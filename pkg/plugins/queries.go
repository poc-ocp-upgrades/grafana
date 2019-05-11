package plugins

import (
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
)

func GetPluginSettings(orgId int64) (map[string]*m.PluginSettingInfoDTO, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	query := m.GetPluginSettingsQuery{OrgId: orgId}
	if err := bus.Dispatch(&query); err != nil {
		return nil, err
	}
	pluginMap := make(map[string]*m.PluginSettingInfoDTO)
	for _, plug := range query.Result {
		pluginMap[plug.PluginId] = plug
	}
	for _, pluginDef := range Plugins {
		if _, ok := pluginMap[pluginDef.Id]; ok {
			continue
		}
		opt := &m.PluginSettingInfoDTO{PluginId: pluginDef.Id, OrgId: orgId, Enabled: true}
		if pluginDef.Type == PluginTypeApp {
			opt.Enabled = false
		}
		if pluginDef.IncludedInAppId != "" {
			opt.Enabled = false
			if appSettings, ok := pluginMap[pluginDef.IncludedInAppId]; ok {
				opt.Enabled = appSettings.Enabled
			}
		}
		pluginMap[pluginDef.Id] = opt
	}
	return pluginMap, nil
}
func GetEnabledPlugins(orgId int64) (*EnabledPlugins, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	enabledPlugins := NewEnabledPlugins()
	pluginSettingMap, err := GetPluginSettings(orgId)
	if err != nil {
		return nil, err
	}
	isPluginEnabled := func(pluginId string) bool {
		_, ok := pluginSettingMap[pluginId]
		return ok
	}
	for pluginId, app := range Apps {
		if b, ok := pluginSettingMap[pluginId]; ok {
			app.Pinned = b.Pinned
			enabledPlugins.Apps = append(enabledPlugins.Apps, app)
		}
	}
	for dsId, ds := range DataSources {
		if isPluginEnabled(ds.Id) {
			enabledPlugins.DataSources[dsId] = ds
		}
	}
	for _, panel := range Panels {
		if isPluginEnabled(panel.Id) {
			enabledPlugins.Panels = append(enabledPlugins.Panels, panel)
		}
	}
	return &enabledPlugins, nil
}
