package notifiers

import (
	"context"
	"time"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
)

const (
	triggMetrString = "Triggered metrics:\n\n"
)

type NotifierBase struct {
	Name			string
	Type			string
	Id			int64
	IsDeault		bool
	UploadImage		bool
	SendReminder		bool
	DisableResolveMessage	bool
	Frequency		time.Duration
	log			log.Logger
}

func NewNotifierBase(model *models.AlertNotification) NotifierBase {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	uploadImage := true
	value, exist := model.Settings.CheckGet("uploadImage")
	if exist {
		uploadImage = value.MustBool()
	}
	return NotifierBase{Id: model.Id, Name: model.Name, IsDeault: model.IsDefault, Type: model.Type, UploadImage: uploadImage, SendReminder: model.SendReminder, DisableResolveMessage: model.DisableResolveMessage, Frequency: model.Frequency, log: log.New("alerting.notifier." + model.Name)}
}
func (n *NotifierBase) ShouldNotify(ctx context.Context, context *alerting.EvalContext, notiferState *models.AlertNotificationState) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if context.PrevAlertState == context.Rule.State && !n.SendReminder {
		return false
	}
	if context.PrevAlertState == context.Rule.State && n.SendReminder {
		lastNotify := time.Unix(notiferState.UpdatedAt, 0)
		if notiferState.UpdatedAt != 0 && lastNotify.Add(n.Frequency).After(time.Now()) {
			return false
		}
		if context.Rule.State == models.AlertStateOK || context.Rule.State == models.AlertStatePending {
			return false
		}
	}
	if context.PrevAlertState == models.AlertStateUnknown && context.Rule.State == models.AlertStateOK {
		return false
	}
	if context.PrevAlertState == models.AlertStateUnknown && context.Rule.State == models.AlertStatePending {
		return false
	}
	if context.PrevAlertState == models.AlertStatePending && context.Rule.State == models.AlertStateOK {
		return false
	}
	if context.PrevAlertState == models.AlertStateOK && context.Rule.State == models.AlertStatePending {
		return false
	}
	if notiferState.State == models.AlertNotificationStatePending {
		lastUpdated := time.Unix(notiferState.UpdatedAt, 0)
		if lastUpdated.Add(1 * time.Minute).After(time.Now()) {
			return false
		}
	}
	if context.Rule.State == models.AlertStateOK && n.DisableResolveMessage {
		return false
	}
	return true
}
func (n *NotifierBase) GetType() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return n.Type
}
func (n *NotifierBase) NeedsImage() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return n.UploadImage
}
func (n *NotifierBase) GetNotifierId() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return n.Id
}
func (n *NotifierBase) GetIsDefault() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return n.IsDeault
}
func (n *NotifierBase) GetSendReminder() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return n.SendReminder
}
func (n *NotifierBase) GetDisableResolveMessage() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return n.DisableResolveMessage
}
func (n *NotifierBase) GetFrequency() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return n.Frequency
}
