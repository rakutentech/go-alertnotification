package alertnotification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// MsTeam is MessageCard for Team notification
type MsTeam struct {
	Type       string          `json:"@type"`
	Context    string          `json:"@context"`
	Summary    string          `json:"summary"`
	ThemeColor string          `json:"themeColor"`
	Title      string          `json:"title"`
	Sections   []SectionStruct `json:"sections"`
}

// SectionStruct is sub-struct of MsTeam
type SectionStruct struct {
	ActivityTitle    string       `json:"activityTitle"`
	ActivitySubtitle string       `json:"activitySubtitle"`
	ActivityImage    string       `json:"activityImage"`
	Facts            []FactStruct `json:"facts"`
}

// FactStruct is sub-struct of SectionStruct
type FactStruct struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewMsTeam is used to create MsTeam
func NewMsTeam(err error, expandos *Expandos) MsTeam {
	title := os.Getenv("ALERT_CARD_SUBJECT")
	summary := os.Getenv("MS_TEAMS_CARD_SUBJECT")
	errMsg := fmt.Sprintf("%+v", err)
	// apply expandos on card
	if expandos != nil {
		if expandos.MsTeamsAlertCardSubject != "" {
			title = expandos.MsTeamsAlertCardSubject
		}
		if expandos.MsTeamsCardSubject != "" {
			summary = expandos.MsTeamsCardSubject
		}
		if expandos.MsTeamsError != "" {
			errMsg = expandos.MsTeamsError
		}
	}

	notificationCard := MsTeam{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		Summary:    summary,
		ThemeColor: os.Getenv("ALERT_THEME_COLOR"),
		Title:      title,
		Sections: []SectionStruct{
			SectionStruct{
				ActivityTitle:    summary,
				ActivitySubtitle: fmt.Sprintf("error has occured on %v", os.Getenv("APP_NAME")),
				ActivityImage:    "",
				Facts: []FactStruct{
					FactStruct{
						Name:  "Environment:",
						Value: os.Getenv("APP_ENV"),
					},
					FactStruct{
						Name:  "ERROR",
						Value: errMsg,
					},
				},
			},
		},
	}
	return notificationCard
}

// Send is implementation of interface AlertNotification's Send()
func (card *MsTeam) Send() (err error) {
	requestBody, err := json.Marshal(card)
	if err != nil {
		return err
	}

	var client http.Client
	timeout := time.Duration(5 * time.Second)
	proxyURL := os.Getenv("MS_TEAMS_PROXY_URL")

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return err
		}
		transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
		client = http.Client{
			Transport: transport,
			Timeout:   timeout,
		}
	} else {
		client = http.Client{
			Timeout: timeout,
		}
	}

	wb := os.Getenv("MS_TEAMS_WEBHOOK")
	if len(wb) == 0 {
		return errors.New("cannot send alert to MSTeams.MS_TEAMS_WEBHOOK is not set in the environment. ")
	}
	request, err := http.NewRequest("POST", wb, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json")
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(respBody) != "1" {
		return errors.New("cannot push to MSTeams")
	}
	return
}
