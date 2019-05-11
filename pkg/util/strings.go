package util

import (
	"fmt"
	"math"
	"regexp"
	"time"
)

func StringsFallback2(val1 string, val2 string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return stringsFallback(val1, val2)
}
func StringsFallback3(val1 string, val2 string, val3 string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return stringsFallback(val1, val2, val3)
}
func stringsFallback(vals ...string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
func SplitString(str string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(str) == 0 {
		return []string{}
	}
	return regexp.MustCompile("[, ]+").Split(str, -1)
}
func GetAgeString(t time.Time) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.IsZero() {
		return "?"
	}
	sinceNow := time.Since(t)
	minutes := sinceNow.Minutes()
	years := int(math.Floor(minutes / 525600))
	months := int(math.Floor(minutes / 43800))
	days := int(math.Floor(minutes / 1440))
	hours := int(math.Floor(minutes / 60))
	if years > 0 {
		return fmt.Sprintf("%dy", years)
	}
	if months > 0 {
		return fmt.Sprintf("%dM", months)
	}
	if days > 0 {
		return fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	if int(minutes) > 0 {
		return fmt.Sprintf("%dm", int(minutes))
	}
	return "< 1m"
}
