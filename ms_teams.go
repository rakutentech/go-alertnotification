package alertnotification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// MsTeam is Adaptive Card for Team notification
type MsTeam struct {
	Type        string       `json:"type"`
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	ContentType string      `json:"contentType"`
	ContentURL  *string     `json:"contentUrl"`
	Content     cardContent `json:"content"`
}

type cardContent struct {
	Schema      string        `json:"$schema"`
	Type        string        `json:"type"`
	Version     string        `json:"version"`
	AccentColor string        `json:"accentColor"`
	Body        []interface{} `json:"body"`
	Actions     []action      `json:"actions"`
	MSTeams     msTeams       `json:"msteams"`
}

type textBlock struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	ID     string `json:"id,omitempty"`
	Size   string `json:"size,omitempty"`
	Weight string `json:"weight,omitempty"`
	Color  string `json:"color,omitempty"`
}

type fact struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type factSet struct {
	Type  string `json:"type"`
	Facts []fact `json:"facts"`
	ID    string `json:"id"`
}

type codeBlock struct {
	Type        string `json:"type"`
	CodeSnippet string `json:"codeSnippet"`
	FontType    string `json:"fontType"`
	Wrap        bool   `json:"wrap"`
}

type action struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type msTeams struct {
	Width string `json:"width"`
}

// NewMsTeam is used to create MsTeam
func NewMsTeam(err error, expandos *Expandos) MsTeam {
	title := os.Getenv("ALERT_CARD_SUBJECT")
	summary := os.Getenv("MS_TEAMS_CARD_SUBJECT")
	errMsg := fmt.Sprintf("%+v", err)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "hostname_unknown"
	}
	hostname += " " + os.Getenv("APP_NAME")
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

	return MsTeam{
		Type: "message",
		Attachments: []attachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				ContentURL:  nil,
				Content: cardContent{
					Schema:      "http://adaptivecards.io/schemas/adaptive-card.json",
					Type:        "AdaptiveCard",
					Version:     "1.4",
					AccentColor: "bf0000",
					Body: []interface{}{
						textBlock{
							Type:   "TextBlock",
							Text:   title,
							ID:     "title",
							Size:   "large",
							Weight: "bolder",
							Color:  "accent",
						},
						factSet{
							Type: "FactSet",
							Facts: []fact{
								{
									Title: "Title:",
									Value: title,
								},
								{
									Title: "Summary:",
									Value: summary,
								},
								{
									Title: "Hostname:",
									Value: hostname,
								},
							},
							ID: "acFactSet",
						},
						codeBlock{
							Type:        "CodeBlock",
							CodeSnippet: errMsg,
							FontType:    "monospace",
							Wrap:        true,
						},
					},
					MSTeams: msTeams{
						Width: "Full",
					},
				},
			},
		},
	}
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
