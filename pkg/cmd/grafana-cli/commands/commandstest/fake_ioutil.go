package commandstest

import (
	"os"
	"time"
)

type FakeIoUtil struct {
	FakeReadDir		[]os.FileInfo
	FakeIsDirectory	bool
}

func (util *FakeIoUtil) Stat(path string) (os.FileInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &FakeFileInfo{IsDirectory: util.FakeIsDirectory}, nil
}
func (util *FakeIoUtil) RemoveAll(path string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (util *FakeIoUtil) ReadDir(path string) ([]os.FileInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return util.FakeReadDir, nil
}
func (i *FakeIoUtil) ReadFile(filename string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return make([]byte, 0), nil
}

type FakeFileInfo struct{ IsDirectory bool }

func (ffi *FakeFileInfo) IsDir() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ffi.IsDirectory
}
func (ffi FakeFileInfo) Size() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return 1
}
func (ffi FakeFileInfo) Mode() os.FileMode {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return 0777
}
func (ffi FakeFileInfo) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (ffi FakeFileInfo) ModTime() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return time.Time{}
}
func (ffi FakeFileInfo) Sys() interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
