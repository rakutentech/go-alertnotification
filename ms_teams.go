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
func NewMsTeam(err error) MsTeam {
	notificationCard := MsTeam{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		Summary:    os.Getenv("MS_TEAMS_CARD_SUBJECT"),
		ThemeColor: os.Getenv("ALERT_THEME_COLOR"),
		Title:      os.Getenv("ALERT_CARD_SUBJECT"),
		Sections: []SectionStruct{
			SectionStruct{
				ActivityTitle:    os.Getenv("MS_TEAMS_CARD_SUBJECT"),
				ActivitySubtitle: fmt.Sprintf("Error has occured on %v", os.Getenv("APP_NAME")),
				ActivityImage:    "",
				Facts: []FactStruct{
					FactStruct{
						Name:  "Environment:",
						Value: os.Getenv("APP_ENV"),
					},
					FactStruct{
						Name:  "ERROR",
						Value: fmt.Sprintf("%+v", err),
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
		return errors.New("Cannot sent alert to MsTeams.MS_TEAMS_WEBHOOK is not set in the environment. ")
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
		return errors.New("Cannot push to MsTeams")
	}
	return
}
