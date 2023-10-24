# golang-alertnotification

This library supports sending throttled alerts as email and as message card to Ms Teams channel.

## Usage

```bash
    go install github.com/rakutentech/go-alertnotification@latest
```

## Configurations

* This package use golang env variables as settings.

### General Configs


| No  | Environment Variable | default | Explanation                                                    |
| :-- | :------------------- | :------ | :------------------------------------------------------------- |
| 1   | APP_ENV              | ""      | application envinronment to be appeared in email/teams message |
| 2   | APP_NAME             | ""      | application name to be appeared in email/teams message         |
| 3   |                      |         |                                                                |

### Email Configs

| No  | Environment Variable | default | Explanation                                                                 |
| :-- | :------------------- | :------ | :-------------------------------------------------------------------------- |
| 1   | EMAIL_ALERT_ENABLED  | false   | change to "true" to enable                                                  |
| 2   | EMAIL_SENDER         | ""      | *require sender email address                                               |
| 3   | EMAIL_RECEIVERS      | ""      | *require receiver email addresses. Eg. test1@example.com, test2@example.com |
| 4   | SMTP_HOST            | ""      | SMTP server hostname                                                        |
| 5   | SMTP_PORT            | ""      | SMTP server port                                                            |
| 6   | EMAIL_USERNAME       | ""      | SMTP username                                                               |
| 7   | EMAIL_PASSWORD       | ""      | SMTP user's passord                                                         |

### Ms Teams Configs

| No  | Environment Variable   | default | Explanation                 |
| :-- | :--------------------- | :------ | :-------------------------- |
| 1   | MS_TEAMS_ALERT_ENABLED | false   | change to "true" to enable  |
| 2   | MS_TEAMS_CARD_SUBJECT  | ""      | Ms teams card subject       |
| 3   | ALERT_CARD_SUBJECT     | ""      | Alert MessageCard subject   |
| 4   | ALERT_THEME_COLOR      | ""      | Themes color                |
| 5   | MS_TEAMS_WEBHOOK       | ""      | *require Ms Teams webhook.  |
| 6   | MS_TEAMS_PROXY_URL     | ""      | Work behind corporate proxy |

### Throttling Configs

| No  | Environment Variable   | default                                    | Explanation                   |
| :-- | :--------------------- | :----------------------------------------- | :---------------------------- |
| 1   | THROTTLE_DURATION      | 7                                          | throttling duration in minute |
| 2   | THROTTLE_DISKCACHE_DIR | /tmp/cache/{APP_NAME}_throttler_disk_cache | disk location for throttling  |
| 3   | THROTTLE_ENABLED       | "true"                                     | Disable all together          |

* Reference for using message card :
<https://docs.microsoft.com/en-us/microsoftteams/platform/concepts/cards/cards-reference>
<https://www.lee-ford.co.uk/send-message-cards-with-microsoft-teams/>

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