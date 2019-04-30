package session

import (
	"math/rand"
	"time"
	ms "github.com/go-macaron/session"
	_ "github.com/go-macaron/session/memcache"
	_ "github.com/go-macaron/session/postgres"
	_ "github.com/go-macaron/session/redis"
	"github.com/grafana/grafana/pkg/log"
	"gopkg.in/macaron.v1"
)

const (
	SESS_KEY_USERID		= "uid"
	SESS_KEY_OAUTH_STATE	= "state"
	SESS_KEY_APIKEY		= "apikey_id"
	SESS_KEY_LASTLDAPSYNC	= "last_ldap_sync"
)

var sessionManager *ms.Manager
var sessionOptions *ms.Options
var StartSessionGC func()
var GetSessionCount func() int
var sessionLogger = log.New("session")
var sessionConnMaxLifetime int64

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	StartSessionGC = func() {
		sessionManager.GC()
		sessionLogger.Debug("Session GC")
		time.AfterFunc(time.Duration(sessionOptions.Gclifetime)*time.Second, StartSessionGC)
	}
	GetSessionCount = func() int {
		return sessionManager.Count()
	}
}
func Init(options *ms.Options, connMaxLifetime int64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	sessionOptions = prepareOptions(options)
	sessionConnMaxLifetime = connMaxLifetime
	sessionManager, err = ms.NewManager(options.Provider, *options)
	if err != nil {
		panic(err)
	}
	rndSeconds := 10 + rand.Int63n(180)
	time.AfterFunc(time.Duration(rndSeconds)*time.Second, StartSessionGC)
}
func prepareOptions(opt *ms.Options) *ms.Options {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(opt.Provider) == 0 {
		opt.Provider = "memory"
	}
	if len(opt.ProviderConfig) == 0 {
		opt.ProviderConfig = "data/sessions"
	}
	if len(opt.CookieName) == 0 {
		opt.CookieName = "grafana_sess"
	}
	if len(opt.CookiePath) == 0 {
		opt.CookiePath = "/"
	}
	if opt.Gclifetime == 0 {
		opt.Gclifetime = 3600
	}
	if opt.Maxlifetime == 0 {
		opt.Maxlifetime = opt.Gclifetime
	}
	if opt.IDLength == 0 {
		opt.IDLength = 16
	}
	return opt
}
func GetSession() SessionStore {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &SessionWrapper{manager: sessionManager}
}

type SessionStore interface {
	Set(interface{}, interface{}) error
	Get(interface{}) interface{}
	Delete(interface{}) interface{}
	ID() string
	Release() error
	Destory(*macaron.Context) error
	Start(*macaron.Context) error
	RegenerateId(*macaron.Context) error
}
type SessionWrapper struct {
	session	ms.RawStore
	manager	*ms.Manager
}

func (s *SessionWrapper) Start(c *macaron.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer func() error {
		if err := recover(); err != nil {
			var retryErr error
			s.session, retryErr = s.manager.Start(c)
			return retryErr
		}
		return nil
	}()
	var err error
	s.session, err = s.manager.Start(c)
	return err
}
func (s *SessionWrapper) RegenerateId(c *macaron.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	s.session, err = s.manager.RegenerateId(c)
	return err
}
func (s *SessionWrapper) Set(k interface{}, v interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.session != nil {
		return s.session.Set(k, v)
	}
	return nil
}
func (s *SessionWrapper) Get(k interface{}) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.session != nil {
		return s.session.Get(k)
	}
	return nil
}
func (s *SessionWrapper) Delete(k interface{}) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.session != nil {
		return s.session.Delete(k)
	}
	return nil
}
func (s *SessionWrapper) ID() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.session != nil {
		return s.session.ID()
	}
	return ""
}
func (s *SessionWrapper) Release() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.session != nil {
		return s.session.Release()
	}
	return nil
}
func (s *SessionWrapper) Destory(c *macaron.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.session != nil {
		if err := s.manager.Destory(c); err != nil {
			return err
		}
		s.session = nil
	}
	return nil
}
