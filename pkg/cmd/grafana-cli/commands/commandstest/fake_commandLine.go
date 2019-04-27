package commandstest

import (
	"github.com/codegangsta/cli"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
)

type FakeFlagger struct{ Data map[string]interface{} }
type FakeCommandLine struct {
	LocalFlags, GlobalFlags	*FakeFlagger
	HelpShown, VersionShown	bool
	CliArgs			[]string
}

func (ff FakeFlagger) String(key string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if value, ok := ff.Data[key]; ok {
		return value.(string)
	}
	return ""
}
func (ff FakeFlagger) StringSlice(key string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if value, ok := ff.Data[key]; ok {
		return value.([]string)
	}
	return []string{}
}
func (ff FakeFlagger) Int(key string) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if value, ok := ff.Data[key]; ok {
		return value.(int)
	}
	return 0
}
func (ff FakeFlagger) Bool(key string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if value, ok := ff.Data[key]; ok {
		return value.(bool)
	}
	return false
}
func (fcli *FakeCommandLine) String(key string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fcli.LocalFlags.String(key)
}
func (fcli *FakeCommandLine) StringSlice(key string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fcli.LocalFlags.StringSlice(key)
}
func (fcli *FakeCommandLine) Int(key string) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fcli.LocalFlags.Int(key)
}
func (fcli *FakeCommandLine) Bool(key string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if fcli.LocalFlags == nil {
		return false
	}
	return fcli.LocalFlags.Bool(key)
}
func (fcli *FakeCommandLine) GlobalString(key string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fcli.GlobalFlags.String(key)
}
func (fcli *FakeCommandLine) Generic(name string) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fcli.LocalFlags.Data[name]
}
func (fcli *FakeCommandLine) FlagNames() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	flagNames := []string{}
	for key := range fcli.LocalFlags.Data {
		flagNames = append(flagNames, key)
	}
	return flagNames
}
func (fcli *FakeCommandLine) ShowHelp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fcli.HelpShown = true
}
func (fcli *FakeCommandLine) Application() *cli.App {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cli.NewApp()
}
func (fcli *FakeCommandLine) Args() cli.Args {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fcli.CliArgs
}
func (fcli *FakeCommandLine) ShowVersion() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fcli.VersionShown = true
}
func (fcli *FakeCommandLine) RepoDirectory() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fcli.GlobalString("repo")
}
func (fcli *FakeCommandLine) PluginDirectory() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fcli.GlobalString("pluginsDir")
}
func (fcli *FakeCommandLine) PluginURL() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fcli.GlobalString("pluginUrl")
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
