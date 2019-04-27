package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"github.com/facebookgo/inject"
	"github.com/grafana/grafana/pkg/api"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/login"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/social"
	"golang.org/x/sync/errgroup"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/services/cache"
	"github.com/grafana/grafana/pkg/setting"
	_ "github.com/grafana/grafana/pkg/extensions"
	_ "github.com/grafana/grafana/pkg/metrics"
	_ "github.com/grafana/grafana/pkg/plugins"
	_ "github.com/grafana/grafana/pkg/services/alerting"
	_ "github.com/grafana/grafana/pkg/services/cleanup"
	_ "github.com/grafana/grafana/pkg/services/notifications"
	_ "github.com/grafana/grafana/pkg/services/provisioning"
	_ "github.com/grafana/grafana/pkg/services/rendering"
	_ "github.com/grafana/grafana/pkg/services/search"
	_ "github.com/grafana/grafana/pkg/services/sqlstore"
	_ "github.com/grafana/grafana/pkg/tracing"
)

func NewGrafanaServer() *GrafanaServerImpl {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)
	return &GrafanaServerImpl{context: childCtx, shutdownFn: shutdownFn, childRoutines: childRoutines, log: log.New("server"), cfg: setting.NewCfg()}
}

type GrafanaServerImpl struct {
	context			context.Context
	shutdownFn		context.CancelFunc
	childRoutines		*errgroup.Group
	log			log.Logger
	cfg			*setting.Cfg
	shutdownReason		string
	shutdownInProgress	bool
	RouteRegister		routing.RouteRegister	`inject:""`
	HttpServer		*api.HTTPServer		`inject:""`
}

func (g *GrafanaServerImpl) Run() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g.loadConfiguration()
	g.writePIDFile()
	login.Init()
	social.NewOAuthService()
	serviceGraph := inject.Graph{}
	serviceGraph.Provide(&inject.Object{Value: bus.GetBus()})
	serviceGraph.Provide(&inject.Object{Value: g.cfg})
	serviceGraph.Provide(&inject.Object{Value: routing.NewRouteRegister(middleware.RequestMetrics, middleware.RequestTracing)})
	serviceGraph.Provide(&inject.Object{Value: cache.New(5*time.Minute, 10*time.Minute)})
	services := registry.GetServices()
	for _, service := range services {
		serviceGraph.Provide(&inject.Object{Value: service.Instance})
	}
	serviceGraph.Provide(&inject.Object{Value: g})
	if err := serviceGraph.Populate(); err != nil {
		return fmt.Errorf("Failed to populate service dependency: %v", err)
	}
	for _, service := range services {
		if registry.IsDisabled(service.Instance) {
			continue
		}
		g.log.Info("Initializing " + service.Name)
		if err := service.Instance.Init(); err != nil {
			return fmt.Errorf("Service init failed: %v", err)
		}
	}
	for _, srv := range services {
		descriptor := srv
		service, ok := srv.Instance.(registry.BackgroundService)
		if !ok {
			continue
		}
		if registry.IsDisabled(descriptor.Instance) {
			continue
		}
		g.childRoutines.Go(func() error {
			if g.shutdownInProgress {
				return nil
			}
			err := service.Run(g.context)
			if err != context.Canceled && err != nil {
				g.log.Error("Stopped "+descriptor.Name, "reason", err)
			} else {
				g.log.Info("Stopped "+descriptor.Name, "reason", err)
			}
			g.shutdownInProgress = true
			return err
		})
	}
	sendSystemdNotification("READY=1")
	return g.childRoutines.Wait()
}
func (g *GrafanaServerImpl) loadConfiguration() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := g.cfg.Load(&setting.CommandLineArgs{Config: *configFile, HomePath: *homePath, Args: flag.Args()})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start grafana. error: %s\n", err.Error())
		os.Exit(1)
	}
	g.log.Info("Starting "+setting.ApplicationName, "version", version, "commit", commit, "branch", buildBranch, "compiled", time.Unix(setting.BuildStamp, 0))
	g.cfg.LogConfigSources()
}
func (g *GrafanaServerImpl) Shutdown(reason string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	g.log.Info("Shutdown started", "reason", reason)
	g.shutdownReason = reason
	g.shutdownInProgress = true
	g.shutdownFn()
	g.childRoutines.Wait()
}
func (g *GrafanaServerImpl) Exit(reason error) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	code := 1
	if reason == context.Canceled && g.shutdownReason != "" {
		reason = fmt.Errorf(g.shutdownReason)
		code = 0
	}
	g.log.Error("Server shutdown", "reason", reason)
	return code
}
func (g *GrafanaServerImpl) writePIDFile() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if *pidFile == "" {
		return
	}
	err := os.MkdirAll(filepath.Dir(*pidFile), 0700)
	if err != nil {
		g.log.Error("Failed to verify pid directory", "error", err)
		os.Exit(1)
	}
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(*pidFile, []byte(pid), 0644); err != nil {
		g.log.Error("Failed to write pidfile", "error", err)
		os.Exit(1)
	}
	g.log.Info("Writing PID file", "path", *pidFile, "pid", pid)
}
func sendSystemdNotification(state string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	notifySocket := os.Getenv("NOTIFY_SOCKET")
	if notifySocket == "" {
		return fmt.Errorf("NOTIFY_SOCKET environment variable empty or unset.")
	}
	socketAddr := &net.UnixAddr{Name: notifySocket, Net: "unixgram"}
	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(state))
	conn.Close()
	return err
}
