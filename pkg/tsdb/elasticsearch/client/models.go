package es

import (
	"encoding/json"
	"github.com/grafana/grafana/pkg/tsdb"
)

type SearchRequest struct {
	Index		string
	Interval	tsdb.Interval
	Size		int
	Sort		map[string]interface{}
	Query		*Query
	Aggs		AggArray
	CustomProps	map[string]interface{}
}

func (r *SearchRequest) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root := make(map[string]interface{})
	root["size"] = r.Size
	if len(r.Sort) > 0 {
		root["sort"] = r.Sort
	}
	for key, value := range r.CustomProps {
		root[key] = value
	}
	root["query"] = r.Query
	if len(r.Aggs) > 0 {
		root["aggs"] = r.Aggs
	}
	return json.Marshal(root)
}

type SearchResponseHits struct {
	Hits	[]map[string]interface{}
	Total	int64
}
type SearchResponse struct {
	Error			map[string]interface{}	`json:"error"`
	Aggregations	map[string]interface{}	`json:"aggregations"`
	Hits			*SearchResponseHits		`json:"hits"`
}
type MultiSearchRequest struct{ Requests []*SearchRequest }
type MultiSearchResponse struct {
	Status		int					`json:"status,omitempty"`
	Responses	[]*SearchResponse	`json:"responses"`
}
type Query struct {
	Bool *BoolQuery `json:"bool"`
}
type BoolQuery struct{ Filters []Filter }

func NewBoolQuery() *BoolQuery {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &BoolQuery{Filters: make([]Filter, 0)}
}
func (q *BoolQuery) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root := make(map[string]interface{})
	if len(q.Filters) > 0 {
		if len(q.Filters) == 1 {
			root["filter"] = q.Filters[0]
		} else {
			root["filter"] = q.Filters
		}
	}
	return json.Marshal(root)
}

type Filter interface{}
type QueryStringFilter struct {
	Filter
	Query			string
	AnalyzeWildcard	bool
}

func (f *QueryStringFilter) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root := map[string]interface{}{"query_string": map[string]interface{}{"query": f.Query, "analyze_wildcard": f.AnalyzeWildcard}}
	return json.Marshal(root)
}

type RangeFilter struct {
	Filter
	Key		string
	Gte		string
	Lte		string
	Format	string
}

const DateFormatEpochMS = "epoch_millis"

func (f *RangeFilter) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root := map[string]map[string]map[string]interface{}{"range": {f.Key: {"lte": f.Lte, "gte": f.Gte}}}
	if f.Format != "" {
		root["range"][f.Key]["format"] = f.Format
	}
	return json.Marshal(root)
}

type Aggregation interface{}
type Agg struct {
	Key			string
	Aggregation	*aggContainer
}

func (a *Agg) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root := map[string]interface{}{a.Key: a.Aggregation}
	return json.Marshal(root)
}

type AggArray []*Agg

func (a AggArray) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	aggsMap := make(map[string]Aggregation)
	for _, subAgg := range a {
		aggsMap[subAgg.Key] = subAgg.Aggregation
	}
	return json.Marshal(aggsMap)
}

type aggContainer struct {
	Type		string
	Aggregation	Aggregation
	Aggs		AggArray
}

func (a *aggContainer) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root := map[string]interface{}{a.Type: a.Aggregation}
	if len(a.Aggs) > 0 {
		root["aggs"] = a.Aggs
	}
	return json.Marshal(root)
}

type aggDef struct {
	key			string
	aggregation	*aggContainer
	builders	[]AggBuilder
}

func newAggDef(key string, aggregation *aggContainer) *aggDef {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &aggDef{key: key, aggregation: aggregation, builders: make([]AggBuilder, 0)}
}

type HistogramAgg struct {
	Interval	int		`json:"interval,omitempty"`
	Field		string	`json:"field"`
	MinDocCount	int		`json:"min_doc_count"`
	Missing		*int	`json:"missing,omitempty"`
}
type DateHistogramAgg struct {
	Field			string			`json:"field"`
	Interval		string			`json:"interval,omitempty"`
	MinDocCount		int				`json:"min_doc_count"`
	Missing			*string			`json:"missing,omitempty"`
	ExtendedBounds	*ExtendedBounds	`json:"extended_bounds"`
	Format			string			`json:"format"`
}
type FiltersAggregation struct {
	Filters map[string]interface{} `json:"filters"`
}
type TermsAggregation struct {
	Field		string					`json:"field"`
	Size		int						`json:"size"`
	Order		map[string]interface{}	`json:"order"`
	MinDocCount	*int					`json:"min_doc_count,omitempty"`
	Missing		*string					`json:"missing,omitempty"`
}
type ExtendedBounds struct {
	Min	string	`json:"min"`
	Max	string	`json:"max"`
}
type GeoHashGridAggregation struct {
	Field		string	`json:"field"`
	Precision	int		`json:"precision"`
}
type MetricAggregation struct {
	Field		string
	Settings	map[string]interface{}
}

func (a *MetricAggregation) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root := map[string]interface{}{"field": a.Field}
	for k, v := range a.Settings {
		if k != "" && v != nil {
			root[k] = v
		}
	}
	return json.Marshal(root)
}

type PipelineAggregation struct {
	BucketPath	string
	Settings	map[string]interface{}
}

func (a *PipelineAggregation) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root := map[string]interface{}{"buckets_path": a.BucketPath}
	for k, v := range a.Settings {
		if k != "" && v != nil {
			root[k] = v
		}
	}
	return json.Marshal(root)
}
