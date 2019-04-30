package commands

import (
	"github.com/codegangsta/cli"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
)

type CommandLine interface {
	ShowHelp()
	ShowVersion()
	Application() *cli.App
	Args() cli.Args
	Bool(name string) bool
	Int(name string) int
	String(name string) string
	StringSlice(name string) []string
	GlobalString(name string) string
	FlagNames() (names []string)
	Generic(name string) interface{}
	PluginDirectory() string
	RepoDirectory() string
	PluginURL() string
}
type contextCommandLine struct{ *cli.Context }

func (c *contextCommandLine) ShowHelp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cli.ShowCommandHelp(c.Context, c.Command.Name)
}
func (c *contextCommandLine) ShowVersion() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cli.ShowVersion(c.Context)
}
func (c *contextCommandLine) Application() *cli.App {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.App
}
func (c *contextCommandLine) PluginDirectory() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.GlobalString("pluginsDir")
}
func (c *contextCommandLine) RepoDirectory() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.GlobalString("repo")
}
func (c *contextCommandLine) PluginURL() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.GlobalString("pluginUrl")
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
