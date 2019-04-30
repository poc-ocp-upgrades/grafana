package imguploader

import (
	"context"
	"path"
	"path/filepath"
	"github.com/grafana/grafana/pkg/setting"
)

type LocalUploader struct{}

func (u *LocalUploader) Upload(ctx context.Context, imageOnDiskPath string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	filename := filepath.Base(imageOnDiskPath)
	image_url := setting.ToAbsUrl(path.Join("public/img/attachments", filename))
	return image_url, nil
}
func NewLocalImageUploader() (*LocalUploader, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &LocalUploader{}, nil
}
