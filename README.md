# golang-alertnotification

This library supports sending alert as email and as message card to Ms Teams' channel. It is built with `go version go1.12.9`

[![CircleCI](https://circleci.com/gh/rakutentech/go-alertnotification/tree/master.svg?style=svg)](https://circleci.com/gh/rakutentech/go-alertnotification/tree/master)

## Installing

### *go get

```bash
    go get -u github.com/rakutentech/go-alertnotification
```

## Configurations

* This package use golang environmrnt variable as setting. It uses `os.Getenv` to get the configuration values. You can use any enivronment setting package. However, in the unittest and example example use `godotenv` from `https://github.com/joho/godotenv`.

### General Configs

|No   |Environment Variable    |default   |Explanation   |
|---|---|---|---|
|1   |APP_ENV   | "" | application envinronment to be appeared in email/teams message     |
|2   |APP_NAME   | ""  | application name to be appeared in email/teams message   |
|3   |   |   |   |

### Email Configs

|No   |Environment Variable    |default       |Explanation   |
|---|---|---|---|
|1   |EMAIL_ALERT_ENABLED   |false    |change to "true" to enable   |
|2   |EMAIL_SENDER   |""   | *require sender email address   |
|3   |EMAIL_RECEIVERS   | ""  | *require receiver email addresses. Multiple address separated by comma. eg. test1@example.com, test2@example.com   |
|4   |SMTP_HOST   |""  | SMTP server hostname   |
|5   |SMTP_PORT   |"" | SMTP server port  |
|6   |EMAIL_USERNAME   |""   |SMTP username   |
|7   |EMAIL_PASSWORD   |""   |SMTP user's passord   |

### Ms Teams Configs

|No   |Environment Variable    |default   |Explanation   |
|---|---|---|---|
|1   |MS_TEAMS_ALERT_ENABLED   |false   | change to "true" to enable    |
|2   |MS_TEAMS_CARD_SUBJECT   |""   | Ms teams card subject  |
|3   |ALERT_CARD_SUBJECT   |""   |Alert MessageCard subject   |
|4   |ALERT_THEME_COLOR   |""   |Themes color   |
|5   |MS_TEAMS_WEBHOOK   |""   |*require Ms Teams webhook. <https://docs.microsoft.com/en-us/microsoftteams/platform/concepts/connectors/connectors-using> |
|6   |MS_TEAMS_PROXY_URL  |""   |Work behind corporate proxy   |

### Throttling Configs

|No   |Environment Variable    |default   |Explanation   |
|---|---|---|---|
|1   |THROTTLE_DURATION   | 7 | throttling duration in minute     |
|2   |THROTTLE_DISKCACHE_DIR   | /tmp/cache/{APP_NAME}_throttler_disk_cache  | disk location for throttling    |
|3   |THROTTLE_ENABLED  | "true"  | Enabled by default to avoid sending too many notification. Set it to "false" to disable. Enable this it will send notification only 1 for the same error within `THROTTLE_DURATION`. Otherwise, it will send each occurence of the error. Recommended to be enable. |

* Reference for using message card :
<https://docs.microsoft.com/en-us/microsoftteams/platform/concepts/cards/cards-reference>
<https://www.lee-ford.co.uk/send-message-cards-with-microsoft-teams/>

## Usage

```golang
 //import
 import n "github.com/rakutentech/go-alertnotification"

 //Create New Alert
 alert := n.NewAlert(err, ignoringErr);

 //Send notification
 alert.Notify();

 // To remove all current throttling
 alert.RemoveCurrentThrotting()

```

## Example

### Add configuration

* Create a `.env` file and add the setting value

```markdown
SMTP_HOST=localhost
SMTP_PORT=25
EMAIL_SENDER=test@example.com
EMAIL_RECEIVERS=recevier.test@exmaple.com
EMAIL_ALERT_ENABLED=true
MS_TEAMS_ALERT_ENABLED=

MS_TEAMS_CARD_SUBJECT=test subject
ALERT_THEME_COLOR=ff5864
ALERT_CARD_SUBJECT=Errror card
MS_TEAMS_CARD_SUBJECT=teams card
APP_ENV=local
APP_NAME=golang
MS_TEAMS_WEBHOOK=Teams webhook
```

### Send alert to both email and teams

```golang
package main;

import (
        "errors"
        "os"
        n "github.com/rakutentech/go-alertnotification"
        "github.com/joho/godotenv"
        )
func setEnv() {
        godotenv.Load()
}


func main() {
    // set env variable
    setEnv();
    err := errors.New("To be alerted error");
    ignoringErr := []error{errors.New("Ignore 001"), errors.New("Ignore 002")};
    alert := n.NewAlert(err, ignoringErr);
    alert.Notify();
}
```
