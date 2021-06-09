package alertnotification

import (
	"fmt"
	"os"
)

// Alert struct for specify the ignoring error and the occuring error
type Alert struct {
	Error            error
	DoNotAlertErrors []error
}

// NewAlert creates Alert sturct instance
func NewAlert(err error, doNotAlertErrors []error) Alert {
	a := Alert{
		Error:            err,
		DoNotAlertErrors: doNotAlertErrors,
	}
	return a
}

// AlertNotification is interface that all send notification function satify including send email
type AlertNotification interface {
	Send() error
}

// DoSendNotification is to send the alert to the specified implemenation of the AlertNoticication interface
func DoSendNotification(alert AlertNotification) error {
	return alert.Send()
}

// Notify send and do throttling when error occur
func (a *Alert) Notify() (err error) {
	if a.shouldAlert() {
		err := a.dispatch()
		fmt.Println(err)
		if err != nil {
			return err
		}
	}
	return
}

// Dispatch sends all notification to all registered chanel
func (a *Alert) dispatch() (err error) {
	if shouldMail() {
		fmt.Println("Send mail....")
		e := NewEmailConfig(a.Error)
		err := e.Send()
		if err != nil {
			return err
		}
	}

	if shouldMsTeams() {
		fmt.Println("SendTeams")
		m := NewMsTeam(a.Error)
		err := m.Send()
		if err != nil {
			return err
		}
	}
	return
}

func (a *Alert) shouldAlert() bool {
	if !a.isThrottlingEnabled() {
		//Always alert when throttling is disabled.
		return true
	}

	if a.isDoNotAlert() {
		return false
	}
	t := NewThrottler()
	return !t.IsThrottled(a.Error)
}

func (a *Alert) isDoNotAlert() bool {
	for _, e := range a.DoNotAlertErrors {
		if e.Error() == a.Error.Error() {
			return true
		}
	}
	return false
}

func shouldMsTeams() bool {
	return os.Getenv("MS_TEAMS_ALERT_ENABLED") == "true"
}

func shouldMail() bool {
	return os.Getenv("EMAIL_ALERT_ENABLED") == "true"
}

func (a *Alert) isThrottlingEnabled() bool {
	return os.Getenv("THROTTLE_ENABLED") != "false"
}

// RemoveCurrentThrotting remove all current throttlings.
func (a *Alert) RemoveCurrentThrotting() error {
	t := NewThrottler()
	return t.CleanThrottlingCache()
}
