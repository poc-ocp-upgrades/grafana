package commands

import (
	"errors"
	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	s "github.com/grafana/grafana/pkg/cmd/grafana-cli/services"
)

func validateVersionInput(c CommandLine) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	arg := c.Args().First()
	if arg == "" {
		return errors.New("please specify plugin to list versions for")
	}
	return nil
}
func listversionsCommand(c CommandLine) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateVersionInput(c); err != nil {
		return err
	}
	pluginToList := c.Args().First()
	plugin, err := s.GetPlugin(pluginToList, c.GlobalString("repo"))
	if err != nil {
		return err
	}
	for _, i := range plugin.Versions {
		logger.Infof("%v\n", i.Version)
	}
	return nil
}
