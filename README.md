# go-alertnotification

This library supports sending throttled alerts as email and as message card to Ms Teams channel.

## Usage

```bash
go install github.com/rakutentech/go-alertnotification@latest
```

## Configurations

* This package use golang env variables as settings.

### General Configs


| Env Variable | default | Description                                                   |
| :----------- | :------ | :------------------------------------------------------------ |
| APP_ENV      |         | application environment to be appeared in email/teams message |
| APP_NAME     |         | application name to be appeared in email/teams message        |


### Email Configs

| Env Variable        | default | Description                                                                     |
| :------------------ | :------ | :------------------------------------------------------------------------------ |
| **EMAIL_SENDER**    |         | **required** sender email address                                               |
| **EMAIL_RECEIVERS** |         | **required** receiver email addresses. Eg. `test1@gmail.com`,`test2@gmail.com`  |
| EMAIL_ALERT_ENABLED | false   | change to "true" to enable                                                      |
| SMTP_HOST           |         | SMTP server hostname                                                            |
| SMTP_PORT           |         | SMTP server port                                                                |
| EMAIL_USERNAME      |         | SMTP username                                                                   |
| EMAIL_PASSWORD      |         | SMTP password                                                                   |

### Ms Teams Configs

| Env Variable           | default | Description                    |
| :--------------------- | :------ | :----------------------------- |
| **MS_TEAMS_WEBHOOK**   |         | **required** Ms Teams webhook. |
| MS_TEAMS_ALERT_ENABLED | false   | change to "true" to enable     |
| MS_TEAMS_CARD_SUBJECT  |         | MS teams card subject          |
| ALERT_CARD_SUBJECT     |         | Alert MessageCard subject      |
| ALERT_THEME_COLOR      |         | Themes color                   |
| MS_TEAMS_PROXY_URL     |         | Work behind corporate proxy    |

### Throttling Configs

| Env Variable           | default                                      | Explanation                    |
| :--------------------- | :------------------------------------------- | :----------------------------- |
| THROTTLE_DURATION      | 7                                            | throttling duration in minutes |
| THROTTLE_GRACE_SECONDS | 0                                            | throttling grace in seconds    |
| THROTTLE_DISKCACHE_DIR | `/tmp/cache/{APP_NAME}_throttler_disk_cache` | disk location for throttling   |
| THROTTLE_ENABLED       | true                                         | Disable all together           |

## Usage

### Simple

```go
 //import
 import n "github.com/rakutentech/go-alertnotification"
 err := errors.New("Alert me")
 ignoringErrs := []error{errors.New("Ignore 001"), errors.New("Ignore 002")};

 //Create New Alert
 alert := n.NewAlert(err, ignoringErrs)
 //Send notification
 alert.Notify()
```

### With customized fields

```go
 import n "github.com/rakutentech/go-alertnotification"

 //Create expandos, can keep the field value as configured by removing that field from expandos
 expandos := &n.Expandos{
        EmailBody:                  "This is the customized email body",
        EmailSubject:               "This is the customized email subject",
        MsTeamsCardSubject:         "This is the customized MS Teams card summary",
        MsTeamsAlertCardSubject:    "This is the customized MS Teams card title",
        MsTeamsError:               "This is the customized MS Teams card error message",
 }

 //Create New Alert
 alert := n.NewAlertWithExpandos(err, ignoringErr, expandos)

 //Send notification
 alert.Notify()

 // To remove all current throttling
 alert.RemoveCurrentThrotting()

```