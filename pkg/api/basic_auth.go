package api

import (
	"crypto/subtle"
	macaron "gopkg.in/macaron.v1"
)

func BasicAuthenticatedRequest(req macaron.Request, expectedUser, expectedPass string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	user, pass, ok := req.BasicAuth()
	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(expectedUser)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(expectedPass)) != 1 {
		return false
	}
	return true
}
