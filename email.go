package alertnotification

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// EmailConfig is email setting struct
type EmailConfig struct {
	Username  string
	Password  string
	Host      string
	Port      string
	Sender    string
	Receivers []string // Can use comma for mutliple email
	ErrorObj  error
}

func getReceivers() []string {
	delimeter := ","
	receivers := os.Getenv("EMAIL_RECEIVERS")
	if len(receivers) == 0 {
		return nil
	}
	return strings.Split(receivers, delimeter)
}

// NewEmailConfig create new EmailConfig struct
func NewEmailConfig(err error) EmailConfig {
	config := EmailConfig{
		Username:  os.Getenv("EMAIL_USERNAME"),
		Password:  os.Getenv("EMAIL_PASSWORD"),
		Host:      os.Getenv("SMTP_HOST"),
		Port:      os.Getenv("SMTP_PORT"),
		Sender:    os.Getenv("EMAIL_SENDER"),
		Receivers: getReceivers(),
		ErrorObj:  err,
	}
	return config
}

// Send Alert email
func (ec *EmailConfig) Send() error {
	fmt.Println("sending email ....")
	var err error
	if ec.Receivers == nil {
		return errors.New("notification receivers are empty")
	}
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	messageDetail := "Error: \r\n" + fmt.Sprintf("%+v", ec.ErrorObj)
	message := "To: " + strings.Join(ec.Receivers, ", ") + "\r\n" +
		"From: " + ec.Sender + "\r\n" +
		"Subject: " + os.Getenv("EMAIL_SUBJECT") + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(string(messageDetail)))

	if len(strings.TrimSpace(ec.Username)) != 0 {
		stmpAuth := smtp.PlainAuth("", ec.Username, ec.Password, ec.Host)

		err = smtp.SendMail(
			ec.Host+":"+ec.Port,
			stmpAuth,
			ec.Sender,
			ec.Receivers,
			[]byte(message),
		)
		return err
	}
	fmt.Println("Send with localhost. ......")
	conn, err := smtp.Dial(ec.Host + ":" + ec.Port)
	if err != nil {
		return err
	}

	defer conn.Close()
	if err = conn.Mail(r.Replace(ec.Sender)); err != nil {
		return err
	}
	// format receiver email
	for i := range ec.Receivers {
		ec.Receivers[i] = r.Replace(ec.Receivers[i])
		if err = conn.Rcpt(ec.Receivers[i]); err != nil {
			return err
		}
	}

	w, err := conn.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return conn.Quit()

}
