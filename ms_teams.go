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

// MsTeam is AdaptiveCard for Team notification
type MsTeam struct {
	Type    string          `json:"type"`
	Version string          `json:"version"`
	Body    []BodyStruct    `json:"body"`
	Actions []ActionStruct  `json:"actions"`
}

// BodyStruct is sub-struct of MsTeam
type BodyStruct struct {
	Type     string       `json:"type"`
	Text     string       `json:"text"`
	Items    []ItemStruct `json:"items"`
}

// ItemStruct is sub-struct of BodyStruct
type ItemStruct struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Weight string `json:"weight"`
	Size   string `json:"size"`
}

// ActionStruct is sub-struct of MsTeam
type ActionStruct struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
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
		Type:    "AdaptiveCard",
		Version: "1.2",
		Body: []BodyStruct{
			BodyStruct{
				Type: "TextBlock",
				Text: title,
				Items: []ItemStruct{
					ItemStruct{
						Type:  "TextBlock",
						Text:  summary,
						Weight: "Bolder",
						Size:   "Medium",
					},
					ItemStruct{
						Type:  "TextBlock",
						Text:  fmt.Sprintf("Error: %s", errMsg),
						Weight: "Lighter",
						Size:   "Small",
					},
				},
			},
		},
		Actions: []ActionStruct{
			ActionStruct{
				Type:  "Action.OpenUrl",
				Title: "View Details",
				URL:   os.Getenv("MS_TEAMS_WEBHOOK"),
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
