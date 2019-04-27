package es

import (
	"strings"
	"github.com/grafana/grafana/pkg/tsdb"
)

type SearchRequestBuilder struct {
	version		int
	interval	tsdb.Interval
	index		string
	size		int
	sort		map[string]interface{}
	queryBuilder	*QueryBuilder
	aggBuilders	[]AggBuilder
	customProps	map[string]interface{}
}

func NewSearchRequestBuilder(version int, interval tsdb.Interval) *SearchRequestBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	builder := &SearchRequestBuilder{version: version, interval: interval, sort: make(map[string]interface{}), customProps: make(map[string]interface{}), aggBuilders: make([]AggBuilder, 0)}
	return builder
}
func (b *SearchRequestBuilder) Build() (*SearchRequest, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sr := SearchRequest{Index: b.index, Interval: b.interval, Size: b.size, Sort: b.sort, CustomProps: b.customProps}
	if b.queryBuilder != nil {
		q, err := b.queryBuilder.Build()
		if err != nil {
			return nil, err
		}
		sr.Query = q
	}
	if len(b.aggBuilders) > 0 {
		sr.Aggs = make(AggArray, 0)
		for _, ab := range b.aggBuilders {
			aggArray, err := ab.Build()
			if err != nil {
				return nil, err
			}
			sr.Aggs = append(sr.Aggs, aggArray...)
		}
	}
	return &sr, nil
}
func (b *SearchRequestBuilder) Size(size int) *SearchRequestBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b.size = size
	return b
}
func (b *SearchRequestBuilder) SortDesc(field, unmappedType string) *SearchRequestBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	props := map[string]string{"order": "desc"}
	if unmappedType != "" {
		props["unmapped_type"] = unmappedType
	}
	b.sort[field] = props
	return b
}
func (b *SearchRequestBuilder) AddDocValueField(field string) *SearchRequestBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if b.version < 5 {
		b.customProps["fields"] = []string{"*", "_source"}
	}
	b.customProps["script_fields"] = make(map[string]interface{})
	if b.version < 5 {
		b.customProps["fielddata_fields"] = []string{field}
	} else {
		b.customProps["docvalue_fields"] = []string{field}
	}
	return b
}
func (b *SearchRequestBuilder) Query() *QueryBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if b.queryBuilder == nil {
		b.queryBuilder = NewQueryBuilder()
	}
	return b.queryBuilder
}
func (b *SearchRequestBuilder) Agg() AggBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	aggBuilder := newAggBuilder(b.version)
	b.aggBuilders = append(b.aggBuilders, aggBuilder)
	return aggBuilder
}

type MultiSearchRequestBuilder struct {
	version		int
	requestBuilders	[]*SearchRequestBuilder
}

func NewMultiSearchRequestBuilder(version int) *MultiSearchRequestBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &MultiSearchRequestBuilder{version: version}
}
func (m *MultiSearchRequestBuilder) Search(interval tsdb.Interval) *SearchRequestBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b := NewSearchRequestBuilder(m.version, interval)
	m.requestBuilders = append(m.requestBuilders, b)
	return b
}
func (m *MultiSearchRequestBuilder) Build() (*MultiSearchRequest, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	requests := []*SearchRequest{}
	for _, sb := range m.requestBuilders {
		searchRequest, err := sb.Build()
		if err != nil {
			return nil, err
		}
		requests = append(requests, searchRequest)
	}
	return &MultiSearchRequest{Requests: requests}, nil
}

type QueryBuilder struct{ boolQueryBuilder *BoolQueryBuilder }

func NewQueryBuilder() *QueryBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &QueryBuilder{}
}
func (b *QueryBuilder) Build() (*Query, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	q := Query{}
	if b.boolQueryBuilder != nil {
		b, err := b.boolQueryBuilder.Build()
		if err != nil {
			return nil, err
		}
		q.Bool = b
	}
	return &q, nil
}
func (b *QueryBuilder) Bool() *BoolQueryBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if b.boolQueryBuilder == nil {
		b.boolQueryBuilder = NewBoolQueryBuilder()
	}
	return b.boolQueryBuilder
}

type BoolQueryBuilder struct{ filterQueryBuilder *FilterQueryBuilder }

func NewBoolQueryBuilder() *BoolQueryBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &BoolQueryBuilder{}
}
func (b *BoolQueryBuilder) Filter() *FilterQueryBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if b.filterQueryBuilder == nil {
		b.filterQueryBuilder = NewFilterQueryBuilder()
	}
	return b.filterQueryBuilder
}
func (b *BoolQueryBuilder) Build() (*BoolQuery, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	boolQuery := BoolQuery{}
	if b.filterQueryBuilder != nil {
		filters, err := b.filterQueryBuilder.Build()
		if err != nil {
			return nil, err
		}
		boolQuery.Filters = filters
	}
	return &boolQuery, nil
}

type FilterQueryBuilder struct{ filters []Filter }

func NewFilterQueryBuilder() *FilterQueryBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &FilterQueryBuilder{filters: make([]Filter, 0)}
}
func (b *FilterQueryBuilder) Build() ([]Filter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b.filters, nil
}
func (b *FilterQueryBuilder) AddDateRangeFilter(timeField, lte, gte, format string) *FilterQueryBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b.filters = append(b.filters, &RangeFilter{Key: timeField, Lte: lte, Gte: gte, Format: format})
	return b
}
func (b *FilterQueryBuilder) AddQueryStringFilter(querystring string, analyseWildcard bool) *FilterQueryBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(strings.TrimSpace(querystring)) == 0 {
		return b
	}
	b.filters = append(b.filters, &QueryStringFilter{Query: querystring, AnalyzeWildcard: analyseWildcard})
	return b
}

type AggBuilder interface {
	Histogram(key, field string, fn func(a *HistogramAgg, b AggBuilder)) AggBuilder
	DateHistogram(key, field string, fn func(a *DateHistogramAgg, b AggBuilder)) AggBuilder
	Terms(key, field string, fn func(a *TermsAggregation, b AggBuilder)) AggBuilder
	Filters(key string, fn func(a *FiltersAggregation, b AggBuilder)) AggBuilder
	GeoHashGrid(key, field string, fn func(a *GeoHashGridAggregation, b AggBuilder)) AggBuilder
	Metric(key, metricType, field string, fn func(a *MetricAggregation)) AggBuilder
	Pipeline(key, pipelineType, bucketPath string, fn func(a *PipelineAggregation)) AggBuilder
	Build() (AggArray, error)
}
type aggBuilderImpl struct {
	AggBuilder
	aggDefs	[]*aggDef
	version	int
}

func newAggBuilder(version int) *aggBuilderImpl {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &aggBuilderImpl{aggDefs: make([]*aggDef, 0), version: version}
}
func (b *aggBuilderImpl) Build() (AggArray, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	aggs := make(AggArray, 0)
	for _, aggDef := range b.aggDefs {
		agg := &Agg{Key: aggDef.key, Aggregation: aggDef.aggregation}
		for _, cb := range aggDef.builders {
			childAggs, err := cb.Build()
			if err != nil {
				return nil, err
			}
			agg.Aggregation.Aggs = append(agg.Aggregation.Aggs, childAggs...)
		}
		aggs = append(aggs, agg)
	}
	return aggs, nil
}
func (b *aggBuilderImpl) Histogram(key, field string, fn func(a *HistogramAgg, b AggBuilder)) AggBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	innerAgg := &HistogramAgg{Field: field}
	aggDef := newAggDef(key, &aggContainer{Type: "histogram", Aggregation: innerAgg})
	if fn != nil {
		builder := newAggBuilder(b.version)
		aggDef.builders = append(aggDef.builders, builder)
		fn(innerAgg, builder)
	}
	b.aggDefs = append(b.aggDefs, aggDef)
	return b
}
func (b *aggBuilderImpl) DateHistogram(key, field string, fn func(a *DateHistogramAgg, b AggBuilder)) AggBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	innerAgg := &DateHistogramAgg{Field: field}
	aggDef := newAggDef(key, &aggContainer{Type: "date_histogram", Aggregation: innerAgg})
	if fn != nil {
		builder := newAggBuilder(b.version)
		aggDef.builders = append(aggDef.builders, builder)
		fn(innerAgg, builder)
	}
	b.aggDefs = append(b.aggDefs, aggDef)
	return b
}

const termsOrderTerm = "_term"

func (b *aggBuilderImpl) Terms(key, field string, fn func(a *TermsAggregation, b AggBuilder)) AggBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	innerAgg := &TermsAggregation{Field: field, Order: make(map[string]interface{})}
	aggDef := newAggDef(key, &aggContainer{Type: "terms", Aggregation: innerAgg})
	if fn != nil {
		builder := newAggBuilder(b.version)
		aggDef.builders = append(aggDef.builders, builder)
		fn(innerAgg, builder)
	}
	if b.version >= 60 && len(innerAgg.Order) > 0 {
		if orderBy, exists := innerAgg.Order[termsOrderTerm]; exists {
			innerAgg.Order["_key"] = orderBy
			delete(innerAgg.Order, termsOrderTerm)
		}
	}
	b.aggDefs = append(b.aggDefs, aggDef)
	return b
}
func (b *aggBuilderImpl) Filters(key string, fn func(a *FiltersAggregation, b AggBuilder)) AggBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	innerAgg := &FiltersAggregation{Filters: make(map[string]interface{})}
	aggDef := newAggDef(key, &aggContainer{Type: "filters", Aggregation: innerAgg})
	if fn != nil {
		builder := newAggBuilder(b.version)
		aggDef.builders = append(aggDef.builders, builder)
		fn(innerAgg, builder)
	}
	b.aggDefs = append(b.aggDefs, aggDef)
	return b
}
func (b *aggBuilderImpl) GeoHashGrid(key, field string, fn func(a *GeoHashGridAggregation, b AggBuilder)) AggBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	innerAgg := &GeoHashGridAggregation{Field: field, Precision: 5}
	aggDef := newAggDef(key, &aggContainer{Type: "geohash_grid", Aggregation: innerAgg})
	if fn != nil {
		builder := newAggBuilder(b.version)
		aggDef.builders = append(aggDef.builders, builder)
		fn(innerAgg, builder)
	}
	b.aggDefs = append(b.aggDefs, aggDef)
	return b
}
func (b *aggBuilderImpl) Metric(key, metricType, field string, fn func(a *MetricAggregation)) AggBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	innerAgg := &MetricAggregation{Field: field, Settings: make(map[string]interface{})}
	aggDef := newAggDef(key, &aggContainer{Type: metricType, Aggregation: innerAgg})
	if fn != nil {
		fn(innerAgg)
	}
	b.aggDefs = append(b.aggDefs, aggDef)
	return b
}
func (b *aggBuilderImpl) Pipeline(key, pipelineType, bucketPath string, fn func(a *PipelineAggregation)) AggBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	innerAgg := &PipelineAggregation{BucketPath: bucketPath, Settings: make(map[string]interface{})}
	aggDef := newAggDef(key, &aggContainer{Type: pipelineType, Aggregation: innerAgg})
	if fn != nil {
		fn(innerAgg)
	}
	b.aggDefs = append(b.aggDefs, aggDef)
	return b
}
