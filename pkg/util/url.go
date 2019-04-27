package util

import (
	"net/url"
	"strings"
)

type UrlQueryReader struct{ values url.Values }

func NewUrlQueryReader(urlInfo *url.URL) (*UrlQueryReader, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	u, err := url.ParseQuery(urlInfo.RawQuery)
	if err != nil {
		return nil, err
	}
	return &UrlQueryReader{values: u}, nil
}
func (r *UrlQueryReader) Get(name string, def string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	val := r.values[name]
	if len(val) == 0 {
		return def
	}
	return val[0]
}
func JoinUrlFragments(a, b string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	if len(b) == 0 {
		return a
	}
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
