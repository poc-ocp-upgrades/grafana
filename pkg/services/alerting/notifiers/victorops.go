package notifiers

import (
	"time"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
	"github.com/grafana/grafana/pkg/setting"
)

const AlertStateCritical = "CRITICAL"
const AlertStateRecovery = "RECOVERY"

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	alerting.RegisterNotifier(&alerting.NotifierPlugin{Type: "victorops", Name: "VictorOps", Description: "Sends notifications to VictorOps", Factory: NewVictoropsNotifier, OptionsTemplate: `
      <h3 class="page-heading">VictorOps settings</h3>
      <div class="gf-form">
        <span class="gf-form-label width-6">Url</span>
        <input type="text" required class="gf-form-input max-width-30" ng-model="ctrl.model.settings.url" placeholder="VictorOps url"></input>
      </div>
      <div class="gf-form">
        <gf-form-switch
           class="gf-form"
           label="Auto resolve incidents"
           label-class="width-14"
           checked="ctrl.model.settings.autoResolve"
           tooltip="Resolve incidents in VictorOps once the alert goes back to ok.">
        </gf-form-switch>
      </div>
    `})
}
func NewVictoropsNotifier(model *models.AlertNotification) (alerting.Notifier, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	autoResolve := model.Settings.Get("autoResolve").MustBool(true)
	url := model.Settings.Get("url").MustString()
	if url == "" {
		return nil, alerting.ValidationError{Reason: "Could not find victorops url property in settings"}
	}
	return &VictoropsNotifier{NotifierBase: NewNotifierBase(model), URL: url, AutoResolve: autoResolve, log: log.New("alerting.notifier.victorops")}, nil
}

type VictoropsNotifier struct {
	NotifierBase
	URL			string
	AutoResolve	bool
	log			log.Logger
}

func (this *VictoropsNotifier) Notify(evalContext *alerting.EvalContext) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	this.log.Info("Executing victorops notification", "ruleId", evalContext.Rule.Id, "notification", this.Name)
	ruleUrl, err := evalContext.GetRuleUrl()
	if err != nil {
		this.log.Error("Failed get rule link", "error", err)
		return err
	}
	if evalContext.Rule.State == models.AlertStateOK && !this.AutoResolve {
		this.log.Info("Not alerting VictorOps", "state", evalContext.Rule.State, "auto resolve", this.AutoResolve)
		return nil
	}
	messageType := evalContext.Rule.State
	if evalContext.Rule.State == models.AlertStateAlerting {
		messageType = AlertStateCritical
	}
	if evalContext.Rule.State == models.AlertStateOK {
		messageType = AlertStateRecovery
	}
	bodyJSON := simplejson.New()
	bodyJSON.Set("message_type", messageType)
	bodyJSON.Set("entity_id", evalContext.Rule.Name)
	bodyJSON.Set("timestamp", time.Now().Unix())
	bodyJSON.Set("state_start_time", evalContext.StartTime.Unix())
	bodyJSON.Set("state_message", evalContext.Rule.Message)
	bodyJSON.Set("monitoring_tool", "Grafana v"+setting.BuildVersion)
	bodyJSON.Set("alert_url", ruleUrl)
	if evalContext.ImagePublicUrl != "" {
		bodyJSON.Set("image_url", evalContext.ImagePublicUrl)
	}
	data, _ := bodyJSON.MarshalJSON()
	cmd := &models.SendWebhookSync{Url: this.URL, Body: string(data)}
	if err := bus.DispatchCtx(evalContext.Ctx, cmd); err != nil {
		this.log.Error("Failed to send Victorops notification", "error", err, "webhook", this.Name)
		return err
	}
	return nil
}
