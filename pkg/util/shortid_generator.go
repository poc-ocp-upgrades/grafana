package util

import (
	"regexp"
	"github.com/teris-io/shortid"
)

var allowedChars = shortid.DefaultABC
var validUidPattern = regexp.MustCompile(`^[a-zA-Z0-9\-\_]*$`).MatchString

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	gen, _ := shortid.New(1, allowedChars, 1)
	shortid.SetDefault(gen)
}
func IsValidShortUid(uid string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return validUidPattern(uid)
}
func GenerateShortUid() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return shortid.MustGenerate()
}
