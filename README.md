# go-alertnotification

This library supports sending throttled alerts as email and as message card to Ms Teams channel.

## Usage

```bash
go install github.com/rakutentech/go-alertnotification@latest
```

## Configurations

* This package use golang env variables as settings.

### General Configs


| No  | Environment Variable | default | Description                                                   |
| :-- | :------------------- | :------ | :------------------------------------------------------------ |
| 1   | APP_ENV              |         | application environment to be appeared in email/teams message |
| 2   | APP_NAME             |         | application name to be appeared in email/teams message        |


### Email Configs

| No  | Environment Variable | default | Description                                                                     |
| :-- | :------------------- | :------ | :------------------------------------------------------------------------------ |
| 1   | EMAIL_ALERT_ENABLED  | false   | change to "true" to enable                                                      |
| 2   | **EMAIL_SENDER**     |         | **required** sender email address                                               |
| 3   | **EMAIL_RECEIVERS**  |         | **required** receiver email addresses. Eg. `test1@gmail.com`, `test2@gmail.com` |
| 4   | SMTP_HOST            |         | SMTP server hostname                                                            |
| 5   | SMTP_PORT            |         | SMTP server port                                                                |
| 6   | EMAIL_USERNAME       |         | SMTP username                                                                   |
| 7   | EMAIL_PASSWORD       |         | SMTP password                                                                   |

### Ms Teams Configs

| No  | Environment Variable   | default | Description                    |
| :-- | :--------------------- | :------ | :----------------------------- |
| 1   | MS_TEAMS_ALERT_ENABLED | false   | change to "true" to enable     |
| 2   | MS_TEAMS_CARD_SUBJECT  |         | MS teams card subject          |
| 3   | ALERT_CARD_SUBJECT     |         | Alert MessageCard subject      |
| 4   | ALERT_THEME_COLOR      |         | Themes color                   |
| 5   | **MS_TEAMS_WEBHOOK**   |         | **required** Ms Teams webhook. |
| 6   | MS_TEAMS_PROXY_URL     |         | Work behind corporate proxy    |

### Throttling Configs

| No  | Environment Variable   | default                                    | Explanation                    |
| :-- | :--------------------- | :----------------------------------------- | :----------------------------- |
| 1   | THROTTLE_DURATION      | 7                                          | throttling duration in minutes |
| 2   | THROTTLE_GRACE_SECONDS | 0                                          | throttling grace in seconds    |
| 3   | THROTTLE_DISKCACHE_DIR | /tmp/cache/{APP_NAME}_throttler_disk_cache | disk location for throttling   |
| 4   | THROTTLE_ENABLED       | true                                       | Disable all together           |

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