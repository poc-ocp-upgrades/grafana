package middleware

import (
	"testing"
	"time"
	"github.com/grafana/grafana/pkg/login"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/session"
	"github.com/grafana/grafana/pkg/setting"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/macaron.v1"
)

func TestAuthProxyWithLdapEnabled(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	Convey("When calling sync grafana user with ldap user", t, func() {
		setting.LdapEnabled = true
		setting.AuthProxyLdapSyncTtl = 60
		servers := []*login.LdapServerConf{{Host: "127.0.0.1"}}
		login.LdapCfg = login.LdapConfig{Servers: servers}
		mockLdapAuther := mockLdapAuthenticator{}
		login.NewLdapAuthenticator = func(server *login.LdapServerConf) login.ILdapAuther {
			return &mockLdapAuther
		}
		Convey("When user logs in, call SyncUser", func() {
			sess := newMockSession()
			ctx := m.ReqContext{Session: &sess}
			So(sess.Get(session.SESS_KEY_LASTLDAPSYNC), ShouldBeNil)
			syncGrafanaUserWithLdapUser(&m.LoginUserQuery{ReqContext: &ctx, Username: "test"})
			So(mockLdapAuther.syncUserCalled, ShouldBeTrue)
			So(sess.Get(session.SESS_KEY_LASTLDAPSYNC), ShouldBeGreaterThan, 0)
		})
		Convey("When session variable not expired, don't sync and don't change session var", func() {
			sess := newMockSession()
			ctx := m.ReqContext{Session: &sess}
			now := time.Now().Unix()
			sess.Set(session.SESS_KEY_LASTLDAPSYNC, now)
			sess.Set(AUTH_PROXY_SESSION_VAR, "test")
			syncGrafanaUserWithLdapUser(&m.LoginUserQuery{ReqContext: &ctx, Username: "test"})
			So(sess.Get(session.SESS_KEY_LASTLDAPSYNC), ShouldEqual, now)
			So(mockLdapAuther.syncUserCalled, ShouldBeFalse)
		})
		Convey("When lastldapsync is expired, session variable should be updated", func() {
			sess := newMockSession()
			ctx := m.ReqContext{Session: &sess}
			expiredTime := time.Now().Add(time.Duration(-120) * time.Minute).Unix()
			sess.Set(session.SESS_KEY_LASTLDAPSYNC, expiredTime)
			sess.Set(AUTH_PROXY_SESSION_VAR, "test")
			syncGrafanaUserWithLdapUser(&m.LoginUserQuery{ReqContext: &ctx, Username: "test"})
			So(sess.Get(session.SESS_KEY_LASTLDAPSYNC), ShouldBeGreaterThan, expiredTime)
			So(mockLdapAuther.syncUserCalled, ShouldBeTrue)
		})
	})
}

type mockSession struct{ value map[interface{}]interface{} }

func newMockSession() mockSession {
	_logClusterCodePath()
	defer _logClusterCodePath()
	session := mockSession{}
	session.value = make(map[interface{}]interface{})
	return session
}
func (s *mockSession) Start(c *macaron.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *mockSession) Set(k interface{}, v interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.value[k] = v
	return nil
}
func (s *mockSession) Get(k interface{}) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.value[k]
}
func (s *mockSession) Delete(k interface{}) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	delete(s.value, k)
	return nil
}
func (s *mockSession) ID() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (s *mockSession) Release() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *mockSession) Destory(c *macaron.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *mockSession) RegenerateId(c *macaron.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}

type mockLdapAuthenticator struct{ syncUserCalled bool }

func (a *mockLdapAuthenticator) Login(query *m.LoginUserQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (a *mockLdapAuthenticator) SyncUser(query *m.LoginUserQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	a.syncUserCalled = true
	return nil
}
func (a *mockLdapAuthenticator) GetGrafanaUserFor(ctx *m.ReqContext, ldapUser *login.LdapUserInfo) (*m.User, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
