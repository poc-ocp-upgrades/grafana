package tsdb

import (
	"fmt"
	"strings"
	"time"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
)

var (
	defaultRes		int64	= 1500
	defaultMinInterval		= time.Millisecond * 1
	year				= time.Hour * 24 * 365
	day				= time.Hour * 24
)

type Interval struct {
	Text	string
	Value	time.Duration
}
type intervalCalculator struct{ minInterval time.Duration }
type IntervalCalculator interface {
	Calculate(timeRange *TimeRange, minInterval time.Duration) Interval
}
type IntervalOptions struct{ MinInterval time.Duration }

func NewIntervalCalculator(opt *IntervalOptions) *intervalCalculator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if opt == nil {
		opt = &IntervalOptions{}
	}
	calc := &intervalCalculator{}
	if opt.MinInterval == 0 {
		calc.minInterval = defaultMinInterval
	} else {
		calc.minInterval = opt.MinInterval
	}
	return calc
}
func (i *Interval) Milliseconds() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return i.Value.Nanoseconds() / int64(time.Millisecond)
}
func (ic *intervalCalculator) Calculate(timerange *TimeRange, minInterval time.Duration) Interval {
	_logClusterCodePath()
	defer _logClusterCodePath()
	to := timerange.MustGetTo().UnixNano()
	from := timerange.MustGetFrom().UnixNano()
	interval := time.Duration((to - from) / defaultRes)
	if interval < minInterval {
		return Interval{Text: formatDuration(minInterval), Value: minInterval}
	}
	rounded := roundInterval(interval)
	return Interval{Text: formatDuration(rounded), Value: rounded}
}
func GetIntervalFrom(dsInfo *models.DataSource, queryModel *simplejson.Json, defaultInterval time.Duration) (time.Duration, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	interval := queryModel.Get("interval").MustString("")
	if interval == "" && dsInfo.JsonData != nil {
		dsInterval := dsInfo.JsonData.Get("timeInterval").MustString("")
		if dsInterval != "" {
			interval = dsInterval
		}
	}
	if interval == "" {
		return defaultInterval, nil
	}
	interval = strings.Replace(strings.Replace(interval, "<", "", 1), ">", "", 1)
	parsedInterval, err := time.ParseDuration(interval)
	if err != nil {
		return time.Duration(0), err
	}
	return parsedInterval, nil
}
func formatDuration(inter time.Duration) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if inter >= year {
		return fmt.Sprintf("%dy", inter/year)
	}
	if inter >= day {
		return fmt.Sprintf("%dd", inter/day)
	}
	if inter >= time.Hour {
		return fmt.Sprintf("%dh", inter/time.Hour)
	}
	if inter >= time.Minute {
		return fmt.Sprintf("%dm", inter/time.Minute)
	}
	if inter >= time.Second {
		return fmt.Sprintf("%ds", inter/time.Second)
	}
	if inter >= time.Millisecond {
		return fmt.Sprintf("%dms", inter/time.Millisecond)
	}
	return "1ms"
}
func roundInterval(interval time.Duration) time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch true {
	case interval <= 15*time.Millisecond:
		return time.Millisecond * 10
	case interval <= 35*time.Millisecond:
		return time.Millisecond * 20
	case interval <= 75*time.Millisecond:
		return time.Millisecond * 50
	case interval <= 150*time.Millisecond:
		return time.Millisecond * 100
	case interval <= 350*time.Millisecond:
		return time.Millisecond * 200
	case interval <= 750*time.Millisecond:
		return time.Millisecond * 500
	case interval <= 1500*time.Millisecond:
		return time.Millisecond * 1000
	case interval <= 3500*time.Millisecond:
		return time.Millisecond * 2000
	case interval <= 7500*time.Millisecond:
		return time.Millisecond * 5000
	case interval <= 12500*time.Millisecond:
		return time.Millisecond * 10000
	case interval <= 17500*time.Millisecond:
		return time.Millisecond * 15000
	case interval <= 25000*time.Millisecond:
		return time.Millisecond * 20000
	case interval <= 45000*time.Millisecond:
		return time.Millisecond * 30000
	case interval <= 90000*time.Millisecond:
		return time.Millisecond * 60000
	case interval <= 210000*time.Millisecond:
		return time.Millisecond * 120000
	case interval <= 450000*time.Millisecond:
		return time.Millisecond * 300000
	case interval <= 750000*time.Millisecond:
		return time.Millisecond * 600000
	case interval <= 1050000*time.Millisecond:
		return time.Millisecond * 900000
	case interval <= 1500000*time.Millisecond:
		return time.Millisecond * 1200000
	case interval <= 2700000*time.Millisecond:
		return time.Millisecond * 1800000
	case interval <= 5400000*time.Millisecond:
		return time.Millisecond * 3600000
	case interval <= 9000000*time.Millisecond:
		return time.Millisecond * 7200000
	case interval <= 16200000*time.Millisecond:
		return time.Millisecond * 10800000
	case interval <= 32400000*time.Millisecond:
		return time.Millisecond * 21600000
	case interval <= 86400000*time.Millisecond:
		return time.Millisecond * 43200000
	case interval <= 172800000*time.Millisecond:
		return time.Millisecond * 86400000
	case interval <= 604800000*time.Millisecond:
		return time.Millisecond * 86400000
	case interval <= 1814400000*time.Millisecond:
		return time.Millisecond * 604800000
	case interval < 3628800000*time.Millisecond:
		return time.Millisecond * 2592000000
	default:
		return time.Millisecond * 31536000000
	}
}
