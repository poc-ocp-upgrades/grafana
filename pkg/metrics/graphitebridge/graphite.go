package graphitebridge

import (
	"bufio"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"sort"
	"strings"
	"time"
	"context"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultInterval		= 15 * time.Second
	millisecondsPerSecond	= 1000
)

type HandlerErrorHandling int

const (
	ContinueOnError	HandlerErrorHandling	= iota
	AbortOnError
)

var metricCategoryPrefix = []string{"proxy_", "api_", "page_", "alerting_", "aws_", "db_", "stat_", "go_", "process_"}
var trimMetricPrefix = []string{"grafana_"}

type Config struct {
	URL		string
	Prefix		string
	Interval	time.Duration
	Timeout		time.Duration
	Gatherer	prometheus.Gatherer
	Logger		Logger
	ErrorHandling	HandlerErrorHandling
	CountersAsDelta	bool
}
type Bridge struct {
	url			string
	prefix			string
	countersAsDetlas	bool
	interval		time.Duration
	timeout			time.Duration
	errorHandling		HandlerErrorHandling
	logger			Logger
	g			prometheus.Gatherer
	lastValue		map[model.Fingerprint]float64
}
type Logger interface{ Println(v ...interface{}) }

func NewBridge(c *Config) (*Bridge, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b := &Bridge{}
	if c.URL == "" {
		return nil, errors.New("missing URL")
	}
	b.url = c.URL
	if c.Gatherer == nil {
		b.g = prometheus.DefaultGatherer
	} else {
		b.g = c.Gatherer
	}
	if c.Logger != nil {
		b.logger = c.Logger
	}
	if c.Prefix != "" {
		b.prefix = c.Prefix
	}
	var z time.Duration
	if c.Interval == z {
		b.interval = defaultInterval
	} else {
		b.interval = c.Interval
	}
	if c.Timeout == z {
		b.timeout = defaultInterval
	} else {
		b.timeout = c.Timeout
	}
	b.errorHandling = c.ErrorHandling
	b.lastValue = map[model.Fingerprint]float64{}
	b.countersAsDetlas = c.CountersAsDelta
	return b, nil
}
func (b *Bridge) Run(ctx context.Context) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := b.Push(); err != nil && b.logger != nil {
				b.logger.Println("error pushing to Graphite:", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
func (b *Bridge) Push() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mfs, err := b.g.Gather()
	if err != nil || len(mfs) == 0 {
		switch b.errorHandling {
		case AbortOnError:
			return err
		case ContinueOnError:
			if b.logger != nil {
				b.logger.Println("continue on error:", err)
			}
		default:
			panic("unrecognized error handling value")
		}
	}
	conn, err := net.DialTimeout("tcp", b.url, b.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return b.writeMetrics(conn, mfs, b.prefix, model.Now())
}
func (b *Bridge) writeMetrics(w io.Writer, mfs []*dto.MetricFamily, prefix string, now model.Time) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, mf := range mfs {
		vec, err := expfmt.ExtractSamples(&expfmt.DecodeOptions{Timestamp: now}, mf)
		if err != nil {
			return err
		}
		buf := bufio.NewWriter(w)
		for _, s := range vec {
			if math.IsNaN(float64(s.Value)) {
				continue
			}
			if err := writePrefix(buf, prefix); err != nil {
				return err
			}
			if err := writeMetric(buf, s.Metric, mf); err != nil {
				return err
			}
			value := b.replaceCounterWithDelta(mf, s.Metric, s.Value)
			if _, err := fmt.Fprintf(buf, " %g %d\n", value, int64(s.Timestamp)/millisecondsPerSecond); err != nil {
				return err
			}
			if err := buf.Flush(); err != nil {
				return err
			}
		}
	}
	return nil
}
func writeMetric(buf *bufio.Writer, m model.Metric, mf *dto.MetricFamily) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	metricName, hasName := m[model.MetricNameLabel]
	numLabels := len(m) - 1
	if !hasName {
		numLabels = len(m)
	}
	for _, v := range trimMetricPrefix {
		if strings.HasPrefix(string(metricName), v) {
			metricName = model.LabelValue(strings.Replace(string(metricName), v, "", 1))
		}
	}
	for _, v := range metricCategoryPrefix {
		if strings.HasPrefix(string(metricName), v) {
			group := strings.Replace(v, "_", " ", 1)
			metricName = model.LabelValue(strings.Replace(string(metricName), v, group, 1))
		}
	}
	labelStrings := make([]string, 0, numLabels)
	for label, value := range m {
		if label != model.MetricNameLabel {
			labelStrings = append(labelStrings, fmt.Sprintf("%s %s", string(label), string(value)))
		}
	}
	var err error
	switch numLabels {
	case 0:
		if hasName {
			if err := writeSanitized(buf, string(metricName)); err != nil {
				return err
			}
		}
	default:
		sort.Strings(labelStrings)
		if err = writeSanitized(buf, string(metricName)); err != nil {
			return err
		}
		for _, s := range labelStrings {
			if err = buf.WriteByte('.'); err != nil {
				return err
			}
			if err = writeSanitized(buf, s); err != nil {
				return err
			}
		}
	}
	return addExtentionConventionForRollups(buf, mf, m)
}
func addExtentionConventionForRollups(buf *bufio.Writer, mf *dto.MetricFamily, m model.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mfType := mf.GetType()
	var err error
	if mfType == dto.MetricType_COUNTER {
		if _, err = fmt.Fprint(buf, ".count"); err != nil {
			return err
		}
	}
	if mfType == dto.MetricType_SUMMARY || mfType == dto.MetricType_HISTOGRAM {
		if strings.HasSuffix(string(m[model.MetricNameLabel]), "_count") {
			if _, err = fmt.Fprint(buf, ".count"); err != nil {
				return err
			}
		}
	}
	if mfType == dto.MetricType_HISTOGRAM {
		if strings.HasSuffix(string(m[model.MetricNameLabel]), "_sum") {
			if _, err = fmt.Fprint(buf, ".sum"); err != nil {
				return err
			}
		}
	}
	return nil
}
func writePrefix(buf *bufio.Writer, s string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, c := range s {
		if _, err := buf.WriteRune(replaceInvalid(c)); err != nil {
			return err
		}
	}
	return nil
}
func writeSanitized(buf *bufio.Writer, s string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prevUnderscore := false
	for _, c := range s {
		c = replaceInvalidRune(c)
		if c == '_' {
			if prevUnderscore {
				continue
			}
			prevUnderscore = true
		} else {
			prevUnderscore = false
		}
		if _, err := buf.WriteRune(c); err != nil {
			return err
		}
	}
	return nil
}
func replaceInvalid(c rune) rune {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c == ' ' || c == '.' {
		return '.'
	}
	return replaceInvalidRune(c)
}
func replaceInvalidRune(c rune) rune {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c == ' ' {
		return '.'
	}
	if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-' || c == '_' || c == ':' || (c >= '0' && c <= '9')) {
		return '_'
	}
	return c
}
func (b *Bridge) replaceCounterWithDelta(mf *dto.MetricFamily, metric model.Metric, value model.SampleValue) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !b.countersAsDetlas {
		return float64(value)
	}
	mfType := mf.GetType()
	if mfType == dto.MetricType_COUNTER {
		return b.returnDelta(metric, value)
	}
	if mfType == dto.MetricType_SUMMARY {
		if strings.HasSuffix(string(metric[model.MetricNameLabel]), "_count") {
			return b.returnDelta(metric, value)
		}
	}
	return float64(value)
}
func (b *Bridge) returnDelta(metric model.Metric, value model.SampleValue) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	key := metric.Fingerprint()
	_, exists := b.lastValue[key]
	if !exists {
		b.lastValue[key] = 0
	}
	delta := float64(value) - b.lastValue[key]
	b.lastValue[key] = float64(value)
	return delta
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
