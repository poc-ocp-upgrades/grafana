package alerting

import (
	"context"
	"fmt"
	"time"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tlog "github.com/opentracing/opentracing-go/log"
	"github.com/benbjohnson/clock"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/services/rendering"
	"github.com/grafana/grafana/pkg/setting"
	"golang.org/x/sync/errgroup"
)

type AlertingService struct {
	RenderService	rendering.Service	`inject:""`
	execQueue	chan *Job
	ticker		*Ticker
	scheduler	Scheduler
	evalHandler	EvalHandler
	ruleReader	RuleReader
	log		log.Logger
	resultHandler	ResultHandler
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registry.RegisterService(&AlertingService{})
}
func NewEngine() *AlertingService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	e := &AlertingService{}
	e.Init()
	return e
}
func (e *AlertingService) IsDisabled() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return !setting.AlertingEnabled || !setting.ExecuteAlerts
}
func (e *AlertingService) Init() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	e.ticker = NewTicker(time.Now(), time.Second*0, clock.New())
	e.execQueue = make(chan *Job, 1000)
	e.scheduler = NewScheduler()
	e.evalHandler = NewEvalHandler()
	e.ruleReader = NewRuleReader()
	e.log = log.New("alerting.engine")
	e.resultHandler = NewResultHandler(e.RenderService)
	return nil
}
func (e *AlertingService) Run(ctx context.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	alertGroup, ctx := errgroup.WithContext(ctx)
	alertGroup.Go(func() error {
		return e.alertingTicker(ctx)
	})
	alertGroup.Go(func() error {
		return e.runJobDispatcher(ctx)
	})
	err := alertGroup.Wait()
	return err
}
func (e *AlertingService) alertingTicker(grafanaCtx context.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer func() {
		if err := recover(); err != nil {
			e.log.Error("Scheduler Panic: stopping alertingTicker", "error", err, "stack", log.Stack(1))
		}
	}()
	tickIndex := 0
	for {
		select {
		case <-grafanaCtx.Done():
			return grafanaCtx.Err()
		case tick := <-e.ticker.C:
			if tickIndex%10 == 0 {
				e.scheduler.Update(e.ruleReader.Fetch())
			}
			e.scheduler.Tick(tick, e.execQueue)
			tickIndex++
		}
	}
}
func (e *AlertingService) runJobDispatcher(grafanaCtx context.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dispatcherGroup, alertCtx := errgroup.WithContext(grafanaCtx)
	for {
		select {
		case <-grafanaCtx.Done():
			return dispatcherGroup.Wait()
		case job := <-e.execQueue:
			dispatcherGroup.Go(func() error {
				return e.processJobWithRetry(alertCtx, job)
			})
		}
	}
}

var (
	unfinishedWorkTimeout	= time.Second * 5
	alertTimeout		= time.Second * 30
	alertMaxAttempts	= 3
)

func (e *AlertingService) processJobWithRetry(grafanaCtx context.Context, job *Job) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer func() {
		if err := recover(); err != nil {
			e.log.Error("Alert Panic", "error", err, "stack", log.Stack(1))
		}
	}()
	cancelChan := make(chan context.CancelFunc, alertMaxAttempts)
	attemptChan := make(chan int, 1)
	attemptChan <- 1
	job.Running = true
	for {
		select {
		case <-grafanaCtx.Done():
			unfinishedWorkTimer := time.NewTimer(unfinishedWorkTimeout)
			select {
			case <-unfinishedWorkTimer.C:
				return e.endJob(grafanaCtx.Err(), cancelChan, job)
			case <-attemptChan:
				return e.endJob(nil, cancelChan, job)
			}
		case attemptID, more := <-attemptChan:
			if !more {
				return e.endJob(nil, cancelChan, job)
			}
			go e.processJob(attemptID, attemptChan, cancelChan, job)
		}
	}
}
func (e *AlertingService) endJob(err error, cancelChan chan context.CancelFunc, job *Job) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	job.Running = false
	close(cancelChan)
	for cancelFn := range cancelChan {
		cancelFn()
	}
	return err
}
func (e *AlertingService) processJob(attemptID int, attemptChan chan int, cancelChan chan context.CancelFunc, job *Job) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer func() {
		if err := recover(); err != nil {
			e.log.Error("Alert Panic", "error", err, "stack", log.Stack(1))
		}
	}()
	alertCtx, cancelFn := context.WithTimeout(context.Background(), alertTimeout)
	cancelChan <- cancelFn
	span := opentracing.StartSpan("alert execution")
	alertCtx = opentracing.ContextWithSpan(alertCtx, span)
	evalContext := NewEvalContext(alertCtx, job.Rule)
	evalContext.Ctx = alertCtx
	go func() {
		defer func() {
			if err := recover(); err != nil {
				e.log.Error("Alert Panic", "error", err, "stack", log.Stack(1))
				ext.Error.Set(span, true)
				span.LogFields(tlog.Error(fmt.Errorf("%v", err)), tlog.String("message", "failed to execute alert rule. panic was recovered."))
				span.Finish()
				close(attemptChan)
			}
		}()
		e.evalHandler.Eval(evalContext)
		span.SetTag("alertId", evalContext.Rule.Id)
		span.SetTag("dashboardId", evalContext.Rule.DashboardId)
		span.SetTag("firing", evalContext.Firing)
		span.SetTag("nodatapoints", evalContext.NoDataFound)
		span.SetTag("attemptID", attemptID)
		if evalContext.Error != nil {
			ext.Error.Set(span, true)
			span.LogFields(tlog.Error(evalContext.Error), tlog.String("message", "alerting execution attempt failed"))
			if attemptID < alertMaxAttempts {
				span.Finish()
				e.log.Debug("Job Execution attempt triggered retry", "timeMs", evalContext.GetDurationMs(), "alertId", evalContext.Rule.Id, "name", evalContext.Rule.Name, "firing", evalContext.Firing, "attemptID", attemptID)
				attemptChan <- (attemptID + 1)
				return
			}
		}
		evalContext.Rule.State = evalContext.GetNewState()
		e.resultHandler.Handle(evalContext)
		span.Finish()
		e.log.Debug("Job Execution completed", "timeMs", evalContext.GetDurationMs(), "alertId", evalContext.Rule.Id, "name", evalContext.Rule.Name, "firing", evalContext.Firing, "attemptID", attemptID)
		close(attemptChan)
	}()
}
