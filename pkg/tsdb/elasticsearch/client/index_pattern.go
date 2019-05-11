package es

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"github.com/grafana/grafana/pkg/tsdb"
)

const (
	noInterval		= ""
	intervalHourly	= "hourly"
	intervalDaily	= "daily"
	intervalWeekly	= "weekly"
	intervalMonthly	= "monthly"
	intervalYearly	= "yearly"
)

type indexPattern interface {
	GetIndices(timeRange *tsdb.TimeRange) ([]string, error)
}

var newIndexPattern = func(interval string, pattern string) (indexPattern, error) {
	if interval == noInterval {
		return &staticIndexPattern{indexName: pattern}, nil
	}
	return newDynamicIndexPattern(interval, pattern)
}

type staticIndexPattern struct{ indexName string }

func (ip *staticIndexPattern) GetIndices(timeRange *tsdb.TimeRange) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []string{ip.indexName}, nil
}

type intervalGenerator interface {
	Generate(from, to time.Time) []time.Time
}
type dynamicIndexPattern struct {
	interval			string
	pattern				string
	intervalGenerator	intervalGenerator
}

func newDynamicIndexPattern(interval, pattern string) (*dynamicIndexPattern, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var generator intervalGenerator
	switch strings.ToLower(interval) {
	case intervalHourly:
		generator = &hourlyInterval{}
	case intervalDaily:
		generator = &dailyInterval{}
	case intervalWeekly:
		generator = &weeklyInterval{}
	case intervalMonthly:
		generator = &monthlyInterval{}
	case intervalYearly:
		generator = &yearlyInterval{}
	default:
		return nil, fmt.Errorf("unsupported interval '%s'", interval)
	}
	return &dynamicIndexPattern{interval: interval, pattern: pattern, intervalGenerator: generator}, nil
}
func (ip *dynamicIndexPattern) GetIndices(timeRange *tsdb.TimeRange) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	from := timeRange.GetFromAsTimeUTC()
	to := timeRange.GetToAsTimeUTC()
	intervals := ip.intervalGenerator.Generate(from, to)
	indices := make([]string, 0)
	for _, t := range intervals {
		indices = append(indices, formatDate(t, ip.pattern))
	}
	return indices, nil
}

type hourlyInterval struct{}

func (i *hourlyInterval) Generate(from, to time.Time) []time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	intervals := []time.Time{}
	start := time.Date(from.Year(), from.Month(), from.Day(), from.Hour(), 0, 0, 0, time.UTC)
	end := time.Date(to.Year(), to.Month(), to.Day(), to.Hour(), 0, 0, 0, time.UTC)
	intervals = append(intervals, start)
	for start.Before(end) {
		start = start.Add(time.Hour)
		intervals = append(intervals, start)
	}
	return intervals
}

type dailyInterval struct{}

func (i *dailyInterval) Generate(from, to time.Time) []time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	intervals := []time.Time{}
	start := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	intervals = append(intervals, start)
	for start.Before(end) {
		start = start.Add(24 * time.Hour)
		intervals = append(intervals, start)
	}
	return intervals
}

type weeklyInterval struct{}

func (i *weeklyInterval) Generate(from, to time.Time) []time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	intervals := []time.Time{}
	start := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	for start.Weekday() != time.Monday {
		start = start.Add(-24 * time.Hour)
	}
	for end.Weekday() != time.Monday {
		end = end.Add(-24 * time.Hour)
	}
	year, week := start.ISOWeek()
	intervals = append(intervals, start)
	for start.Before(end) {
		start = start.Add(24 * time.Hour)
		nextYear, nextWeek := start.ISOWeek()
		if nextYear != year || nextWeek != week {
			intervals = append(intervals, start)
		}
		year = nextYear
		week = nextWeek
	}
	return intervals
}

type monthlyInterval struct{}

func (i *monthlyInterval) Generate(from, to time.Time) []time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	intervals := []time.Time{}
	start := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC)
	month := start.Month()
	intervals = append(intervals, start)
	for start.Before(end) {
		start = start.Add(24 * time.Hour)
		nextMonth := start.Month()
		if nextMonth != month {
			intervals = append(intervals, start)
		}
		month = nextMonth
	}
	return intervals
}

type yearlyInterval struct{}

func (i *yearlyInterval) Generate(from, to time.Time) []time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	intervals := []time.Time{}
	start := time.Date(from.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(to.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	year := start.Year()
	intervals = append(intervals, start)
	for start.Before(end) {
		start = start.Add(24 * time.Hour)
		nextYear := start.Year()
		if nextYear != year {
			intervals = append(intervals, start)
		}
		year = nextYear
	}
	return intervals
}

var datePatternRegex = regexp.MustCompile("(LT|LL?L?L?|l{1,4}|Mo|MM?M?M?|Do|DDDo|DD?D?D?|ddd?d?|do?|w[o|w]?|W[o|W]?|YYYYY|YYYY|YY|gg(ggg?)?|GG(GGG?)?|e|E|a|A|hh?|HH?|mm?|ss?|SS?S?|X|zz?|ZZ?|Q)")
var datePatternReplacements = map[string]string{"M": "1", "MM": "01", "MMM": "Jan", "MMMM": "January", "D": "2", "DD": "02", "DDD": "<stdDayOfYear>", "DDDD": "<stdDayOfYearZero>", "d": "<stdDayOfWeek>", "dd": "Mon", "ddd": "Mon", "dddd": "Monday", "e": "<stdDayOfWeek>", "E": "<stdDayOfWeekISO>", "w": "<stdWeekOfYear>", "ww": "<stdWeekOfYear>", "W": "<stdWeekOfYear>", "WW": "<stdWeekOfYear>", "YY": "06", "YYYY": "2006", "gg": "<stdIsoYearShort>", "gggg": "<stdIsoYear>", "GG": "<stdIsoYearShort>", "GGGG": "<stdIsoYear>", "Q": "<stdQuarter>", "A": "PM", "a": "pm", "H": "<stdHourNoZero>", "HH": "15", "h": "3", "hh": "03", "m": "4", "mm": "04", "s": "5", "ss": "05", "z": "MST", "zz": "MST", "Z": "Z07:00", "ZZ": "-0700", "X": "<stdUnix>", "LT": "3:04 PM", "L": "01/02/2006", "l": "1/2/2006", "ll": "Jan 2 2006", "lll": "Jan 2 2006 3:04 PM", "llll": "Mon, Jan 2 2006 3:04 PM"}

func formatDate(t time.Time, pattern string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var datePattern string
	base := ""
	ltr := false
	if strings.HasPrefix(pattern, "[") {
		parts := strings.Split(strings.TrimLeft(pattern, "["), "]")
		base = parts[0]
		if len(parts) == 2 {
			datePattern = parts[1]
		} else {
			datePattern = base
			base = ""
		}
		ltr = true
	} else if strings.HasSuffix(pattern, "]") {
		parts := strings.Split(strings.TrimRight(pattern, "]"), "[")
		datePattern = parts[0]
		if len(parts) == 2 {
			base = parts[1]
		} else {
			base = ""
		}
		ltr = false
	}
	formatted := t.Format(patternToLayout(datePattern))
	if strings.Contains(formatted, "<std") {
		isoYear, isoWeek := t.ISOWeek()
		isoYearShort := fmt.Sprintf("%d", isoYear)[2:4]
		formatted = strings.Replace(formatted, "<stdIsoYear>", fmt.Sprintf("%d", isoYear), -1)
		formatted = strings.Replace(formatted, "<stdIsoYearShort>", isoYearShort, -1)
		formatted = strings.Replace(formatted, "<stdWeekOfYear>", fmt.Sprintf("%d", isoWeek), -1)
		formatted = strings.Replace(formatted, "<stdUnix>", fmt.Sprintf("%d", t.Unix()), -1)
		day := t.Weekday()
		dayOfWeekIso := int(day)
		if day == time.Sunday {
			dayOfWeekIso = 7
		}
		formatted = strings.Replace(formatted, "<stdDayOfWeek>", fmt.Sprintf("%d", day), -1)
		formatted = strings.Replace(formatted, "<stdDayOfWeekISO>", fmt.Sprintf("%d", dayOfWeekIso), -1)
		formatted = strings.Replace(formatted, "<stdDayOfYear>", fmt.Sprintf("%d", t.YearDay()), -1)
		quarter := 4
		switch t.Month() {
		case time.January, time.February, time.March:
			quarter = 1
		case time.April, time.May, time.June:
			quarter = 2
		case time.July, time.August, time.September:
			quarter = 3
		}
		formatted = strings.Replace(formatted, "<stdQuarter>", fmt.Sprintf("%d", quarter), -1)
		formatted = strings.Replace(formatted, "<stdHourNoZero>", fmt.Sprintf("%d", t.Hour()), -1)
	}
	if ltr {
		return base + formatted
	}
	return formatted + base
}
func patternToLayout(pattern string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var match [][]string
	if match = datePatternRegex.FindAllStringSubmatch(pattern, -1); match == nil {
		return pattern
	}
	for i := range match {
		if replace, ok := datePatternReplacements[match[i][0]]; ok {
			pattern = strings.Replace(pattern, match[i][0], replace, 1)
		}
	}
	return pattern
}
