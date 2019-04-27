package testdata

import (
	"context"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/tsdb"
)

type TestDataExecutor struct {
	*models.DataSource
	log	log.Logger
}

func NewTestDataExecutor(dsInfo *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TestDataExecutor{DataSource: dsInfo, log: log.New("tsdb.testdata")}, nil
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tsdb.RegisterTsdbQueryEndpoint("testdata", NewTestDataExecutor)
}
func (e *TestDataExecutor) Query(ctx context.Context, dsInfo *models.DataSource, tsdbQuery *tsdb.TsdbQuery) (*tsdb.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := &tsdb.Response{}
	result.Results = make(map[string]*tsdb.QueryResult)
	for _, query := range tsdbQuery.Queries {
		scenarioId := query.Model.Get("scenarioId").MustString("random_walk")
		if scenario, exist := ScenarioRegistry[scenarioId]; exist {
			result.Results[query.RefId] = scenario.Handler(query, tsdbQuery)
			result.Results[query.RefId].RefId = query.RefId
		} else {
			e.log.Error("Scenario not found", "scenarioId", scenarioId)
		}
	}
	return result, nil
}
