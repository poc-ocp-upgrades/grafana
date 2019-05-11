package rendering

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	plugin "github.com/hashicorp/go-plugin"
	pluginModel "github.com/grafana/grafana-plugin-model/go/renderer"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registry.RegisterService(&RenderingService{})
}

type RenderingService struct {
	log				log.Logger
	pluginClient	*plugin.Client
	grpcPlugin		pluginModel.RendererPlugin
	pluginInfo		*plugins.RendererPlugin
	renderAction	renderFunc
	domain			string
	inProgressCount	int
	Cfg				*setting.Cfg	`inject:""`
}

func (rs *RenderingService) Init() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rs.log = log.New("rendering")
	err := os.MkdirAll(rs.Cfg.ImagesDir, 0700)
	if err != nil {
		return err
	}
	if rs.Cfg.RendererUrl != "" {
		u, _ := url.Parse(rs.Cfg.RendererCallbackUrl)
		rs.domain = u.Hostname()
	} else if setting.HttpAddr != setting.DEFAULT_HTTP_ADDR {
		rs.domain = setting.HttpAddr
	} else {
		rs.domain = "localhost"
	}
	return nil
}
func (rs *RenderingService) Run(ctx context.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if rs.Cfg.RendererUrl != "" {
		rs.log.Info("Backend rendering via external http server")
		rs.renderAction = rs.renderViaHttp
		<-ctx.Done()
		return nil
	}
	if plugins.Renderer == nil {
		rs.renderAction = rs.renderViaPhantomJS
		<-ctx.Done()
		return nil
	}
	rs.pluginInfo = plugins.Renderer
	if err := rs.startPlugin(ctx); err != nil {
		return err
	}
	rs.renderAction = rs.renderViaPlugin
	err := rs.watchAndRestartPlugin(ctx)
	if rs.pluginClient != nil {
		rs.log.Debug("Killing renderer plugin process")
		rs.pluginClient.Kill()
	}
	return err
}
func (rs *RenderingService) Render(ctx context.Context, opts Opts) (*RenderResult, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if rs.inProgressCount > opts.ConcurrentLimit {
		return &RenderResult{FilePath: filepath.Join(setting.HomePath, "public/img/rendering_limit.png")}, nil
	}
	defer func() {
		rs.inProgressCount -= 1
	}()
	rs.inProgressCount += 1
	if rs.renderAction != nil {
		return rs.renderAction(ctx, opts)
	} else {
		return nil, fmt.Errorf("No renderer found")
	}
}
func (rs *RenderingService) getFilePathForNewImage() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pngPath, _ := filepath.Abs(filepath.Join(rs.Cfg.ImagesDir, util.GetRandomString(20)))
	return pngPath + ".png"
}
func (rs *RenderingService) getURL(path string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if rs.Cfg.RendererUrl != "" {
		return fmt.Sprintf("%s%s&render=1", rs.Cfg.RendererCallbackUrl, path)
	}
	return fmt.Sprintf("%s://%s:%s/%s&render=1", setting.Protocol, rs.domain, setting.HttpPort, path)
}
func (rs *RenderingService) getRenderKey(orgId, userId int64, orgRole models.RoleType) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return middleware.AddRenderAuthKey(orgId, userId, orgRole)
}
