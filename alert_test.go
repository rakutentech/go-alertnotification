package alertnotification

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func setEnv() {
	if err := godotenv.Load("test.env"); err != nil {
		fmt.Println(err)
	}
}

func Test_shouldMsTeams(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "shouldSend", want: true},
		{name: "shouldNotSend", want: false},
	}

	for _, tt := range tests {
		setEnv()
		switch tt.name {
		case "shouldSend":
			os.Setenv("EMAIL_ALERT_ENABLED", "true")
			os.Setenv("MS_TEAMS_ALERT_ENABLED", "true")
		case "shouldNotSend":
			os.Setenv("MS_TEAMS_ALERT_ENABLED", "")
		}
		t.Run(tt.name, func(t *testing.T) {

			if got := shouldMsTeams(); got != tt.want {
				t.Errorf("shouldMsTeams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shouldMail(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{

		{name: "shouldSend", want: true},
		{name: "shouldNotSend", want: false},
	}

	for _, tt := range tests {
		setEnv()
		switch tt.name {
		case "shouldSend":
			os.Setenv("EMAIL_ALERT_ENABLED", "true")
		case "shouldNotSend":
			os.Setenv("EMAIL_ALERT_ENABLED", "")
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldMail(); got != tt.want {
				t.Errorf("shouldMail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlert_Notify(t *testing.T) {
	expandos := make(map[string]string)
	expandos["body"] = "This is mail body"
	expandos["subject"] = "This is mail subject"
	type fields struct {
		Error            error
		DoNotAlertErrors []error
		Expandos         map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Notify_false",
			fields: fields{
				Error: errors.New("Do not alert"), // Do not alert => no error will occur
				DoNotAlertErrors: []error{
					errors.New("Do not alert"), errors.New("if this error then alert")},
				Expandos: nil,
			},
			wantErr: false,
		},
		{name: "Notify_true",
			fields: fields{
				Error: errors.New("give an alert"), // error occured and try to send email => no mail setting configure => error.
				DoNotAlertErrors: []error{
					errors.New("Do not alert"), errors.New("if this error then alert")},
				Expandos: nil,
			},
			wantErr: true,
		},
		{name: "Expandos",
			fields: fields{
				Error: errors.New("give an alert"), // error occured and try to send email => no mail setting configure => error.
				DoNotAlertErrors: []error{
					errors.New("Do not alert"), errors.New("if this error then alert")},
				Expandos: expandos,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if err := godotenv.Overload("test.env"); err != nil { // Reload Env
			fmt.Println(err)
		}
		t.Run(tt.name, func(t *testing.T) {

			a := &Alert{
				Error:            tt.fields.Error,
				DoNotAlertErrors: tt.fields.DoNotAlertErrors,
				Expandos:         tt.fields.Expandos,
			}
			if err := a.RemoveCurrentThrotting(); err != nil {
				t.Errorf("Alert.Notify() error = %+v", err)
			}
			if err := a.Notify(); (err != nil) != tt.wantErr {
				t.Errorf("Alert.Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlert_shouldAlert(t *testing.T) {
	type fields struct {
		Error            error
		DoNotAlertErrors []error
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "shouldAlert_false",
			fields: fields{
				Error: errors.New("Do not alert"),
				DoNotAlertErrors: []error{
					errors.New("Do not alert"), errors.New("if this error then don't alert")},
			},
			want: false,
		},
		{name: "shouldAlert_true",
			fields: fields{
				Error: errors.New("alert this"),
				DoNotAlertErrors: []error{
					errors.New("do not alert"), errors.New("if this error then don't alert")},
			},
			want: true,
		},
		{name: "shouldAlert_true_disable_throttling",
			fields: fields{
				Error: errors.New("do not alert"),
				DoNotAlertErrors: []error{
					errors.New("do not alert"), errors.New("if this error then don't alert")},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "shouldAlert_true_disable_throttling" {
				os.Setenv("THROTTLE_ENABLED", "false")
			}
			a := &Alert{
				Error:            tt.fields.Error,
				DoNotAlertErrors: tt.fields.DoNotAlertErrors,
			}
			if err := a.RemoveCurrentThrotting(); err != nil {
				t.Errorf("Alert.Notify() error = %+v", err)
			}
			if got := a.shouldAlert(); got != tt.want {
				t.Errorf("Alert.shouldAlert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlert_isDoNotAlert(t *testing.T) {
	type fields struct {
		Error            error
		DoNotAlertErrors []error
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "isDoNotAlert_true",
			fields: fields{
				Error: errors.New("Do not alert"),
				DoNotAlertErrors: []error{
					errors.New("Do not alert"), errors.New("if this error then not alert")},
			},
			want: true,
		},
		{name: "isDoNotAlert_false",
			fields: fields{
				Error: errors.New("give an alert"),
				DoNotAlertErrors: []error{
					errors.New("Do not alert"), errors.New("if this error then do alert")},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Alert{
				Error:            tt.fields.Error,
				DoNotAlertErrors: tt.fields.DoNotAlertErrors,
			}
			if got := a.isDoNotAlert(); got != tt.want {
				t.Errorf("Alert.isDoNotAlert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlert_isThrottlingEnabled(t *testing.T) {
	type fields struct {
		Error            error
		DoNotAlertErrors []error
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "isThrottlingEnabled_true",
			fields: fields{
				Error: errors.New("Do not alert"),
				DoNotAlertErrors: []error{
					errors.New("Do not alert"), errors.New("if this error then not alert")},
			},
			want: true,
		},
		{name: "isThrottlingEnabled_false",
			fields: fields{
				Error: errors.New("give an alert"),
				DoNotAlertErrors: []error{
					errors.New("Do not alert"), errors.New("if this error then do alert")},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("THROTTLE_ENABLED", "true")
			if tt.name == "isThrottlingEnabled_false" {
				os.Setenv("THROTTLE_ENABLED", "false")
			}
			a := &Alert{
				Error:            tt.fields.Error,
				DoNotAlertErrors: tt.fields.DoNotAlertErrors,
			}
			if got := a.isThrottlingEnabled(); got != tt.want {
				t.Errorf("Alert.isThrottlingEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
