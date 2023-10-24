<p align="center">
  <a href="https://github.com/rakutentech/go-alertnotification">
    <img alt="go-alertnotification" src="logo.png" width="360">
  </a>
</p>

<h3 align="center">Send Alert Notifications for Go Errors.<br>Notify when new error arrives.</h3>

<p align="center">
  Throttle notifications to avoid overwhelming your inbox.
  <br>
  Grace period to notify if the same error occurs.
  <br>
  Supports multiple Emails, MS Teams and proxy support.
</p>

## Usage

```bash
go install github.com/rakutentech/go-alertnotification@latest
```

## Configurations

You need to set env variables to configure the alert notification behaviour.

### General Configs


| Env Variable | default | Description                  |
| :----------- | :------ | :--------------------------- |
| APP_ENV      |         | appears in email/teams title |
| APP_NAME     |         | appears in email/teams body  |


### Email Configs

| Env Variable        | default | Description                                                              |
| :------------------ | :------ | :----------------------------------------------------------------------- |
| **EMAIL_SENDER**    |         | **Required** one mail address                                            |
| **EMAIL_RECEIVERS** |         | **Required** multiple addresses. Eg. `test1@gmail.com`,`test2@gmail.com` |
| EMAIL_ALERT_ENABLED | false   |                                                                          |
| SMTP_HOST           |         |                                                                          |
| SMTP_PORT           |         |                                                                          |
| EMAIL_USERNAME      |         |                                                                          |
| EMAIL_PASSWORD      |         |                                                                          |

### Ms Teams Configs

| Env Variable           | default | Description        |
| :--------------------- | :------ | :----------------- |
| **MS_TEAMS_WEBHOOK**   |         | **Required**       |
| MS_TEAMS_ALERT_ENABLED | false   |                    |
| MS_TEAMS_CARD_SUBJECT  |         |                    |
| ALERT_CARD_SUBJECT     |         |                    |
| ALERT_THEME_COLOR      |         |                    |
| MS_TEAMS_PROXY_URL     |         | HTTP proxy, if any |

### Throttling Configs

| Env Variable           | default                         | Explanation                      |
| :--------------------- | :------------------------------ | :------------------------------- |
| THROTTLE_DURATION      | 7                               | throttling duration in (minutes) |
| THROTTLE_GRACE_SECONDS | 0                               | throttling grace in (seconds)    |
| THROTTLE_DISKCACHE_DIR | `/tmp/cache/{APP_NAME}_thro...` | disk location for throttling     |
| THROTTLE_ENABLED       | true                            | Disable all together             |

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