package login

import (
	"testing"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLdapLogin(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	Convey("Login using ldap", t, func() {
		Convey("Given ldap enabled and a server configured", func() {
			setting.LdapEnabled = true
			LdapCfg.Servers = append(LdapCfg.Servers, &LdapServerConf{Host: ""})
			ldapLoginScenario("When login with invalid credentials", func(sc *ldapLoginScenarioContext) {
				sc.withLoginResult(false)
				enabled, err := loginUsingLdap(sc.loginUserQuery)
				Convey("it should return true", func() {
					So(enabled, ShouldBeTrue)
				})
				Convey("it should return invalid credentials error", func() {
					So(err, ShouldEqual, ErrInvalidCredentials)
				})
				Convey("it should call ldap login", func() {
					So(sc.ldapAuthenticatorMock.loginCalled, ShouldBeTrue)
				})
			})
			ldapLoginScenario("When login with valid credentials", func(sc *ldapLoginScenarioContext) {
				sc.withLoginResult(true)
				enabled, err := loginUsingLdap(sc.loginUserQuery)
				Convey("it should return true", func() {
					So(enabled, ShouldBeTrue)
				})
				Convey("it should not return error", func() {
					So(err, ShouldBeNil)
				})
				Convey("it should call ldap login", func() {
					So(sc.ldapAuthenticatorMock.loginCalled, ShouldBeTrue)
				})
			})
		})
		Convey("Given ldap enabled and no server configured", func() {
			setting.LdapEnabled = true
			LdapCfg.Servers = make([]*LdapServerConf, 0)
			ldapLoginScenario("When login", func(sc *ldapLoginScenarioContext) {
				sc.withLoginResult(true)
				enabled, err := loginUsingLdap(sc.loginUserQuery)
				Convey("it should return true", func() {
					So(enabled, ShouldBeTrue)
				})
				Convey("it should return invalid credentials error", func() {
					So(err, ShouldEqual, ErrInvalidCredentials)
				})
				Convey("it should not call ldap login", func() {
					So(sc.ldapAuthenticatorMock.loginCalled, ShouldBeFalse)
				})
			})
		})
		Convey("Given ldap disabled", func() {
			setting.LdapEnabled = false
			ldapLoginScenario("When login", func(sc *ldapLoginScenarioContext) {
				sc.withLoginResult(false)
				enabled, err := loginUsingLdap(&m.LoginUserQuery{Username: "user", Password: "pwd"})
				Convey("it should return false", func() {
					So(enabled, ShouldBeFalse)
				})
				Convey("it should not return error", func() {
					So(err, ShouldBeNil)
				})
				Convey("it should not call ldap login", func() {
					So(sc.ldapAuthenticatorMock.loginCalled, ShouldBeFalse)
				})
			})
		})
	})
}
func mockLdapAuthenticator(valid bool) *mockLdapAuther {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mock := &mockLdapAuther{validLogin: valid}
	NewLdapAuthenticator = func(server *LdapServerConf) ILdapAuther {
		return mock
	}
	return mock
}

type mockLdapAuther struct {
	validLogin	bool
	loginCalled	bool
}

func (a *mockLdapAuther) Login(query *m.LoginUserQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	a.loginCalled = true
	if !a.validLogin {
		return ErrInvalidCredentials
	}
	return nil
}
func (a *mockLdapAuther) SyncUser(query *m.LoginUserQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (a *mockLdapAuther) GetGrafanaUserFor(ctx *m.ReqContext, ldapUser *LdapUserInfo) (*m.User, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}

type ldapLoginScenarioContext struct {
	loginUserQuery		*m.LoginUserQuery
	ldapAuthenticatorMock	*mockLdapAuther
}
type ldapLoginScenarioFunc func(c *ldapLoginScenarioContext)

func ldapLoginScenario(desc string, fn ldapLoginScenarioFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	Convey(desc, func() {
		origNewLdapAuthenticator := NewLdapAuthenticator
		sc := &ldapLoginScenarioContext{loginUserQuery: &m.LoginUserQuery{Username: "user", Password: "pwd", IpAddress: "192.168.1.1:56433"}, ldapAuthenticatorMock: &mockLdapAuther{}}
		defer func() {
			NewLdapAuthenticator = origNewLdapAuthenticator
		}()
		fn(sc)
	})
}
func (sc *ldapLoginScenarioContext) withLoginResult(valid bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sc.ldapAuthenticatorMock = mockLdapAuthenticator(valid)
}
