package plugins

import "encoding/json"

type RendererPlugin struct {
	PluginBase
	Executable	string	`json:"executable,omitempty"`
}

func (r *RendererPlugin) Load(decoder *json.Decoder, pluginDir string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := decoder.Decode(&r); err != nil {
		return err
	}
	if err := r.registerPlugin(pluginDir); err != nil {
		return err
	}
	Renderer = r
	return nil
}
