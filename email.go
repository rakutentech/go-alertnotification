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
	Recievers []string // Can use comma for mutliple email
	ErrorObj  error
}

func getReceivers() []string {
	delimeter := ","
	receivers := os.Getenv("EMAIL_RECIEVERS")
	if len(receivers) == 0 {
		return nil
	}
	receivers = strings.TrimSpace(receivers)
	if strings.HasSuffix(receivers, delimeter) {
		receivers = strings.TrimSuffix(receivers, delimeter)
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
		Recievers: getReceivers(),
		ErrorObj:  err,
	}
	return config
}

// Send Alert email
func (ec *EmailConfig) Send() error {
	fmt.Println("sending email ....")
	var err error
	if ec.Recievers == nil {
		return errors.New("Notification Receivers is empty")
	}
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	message := "To: " + strings.Join(ec.Recievers, ", ") +
		"From: " + ec.Sender + "\r\n" +
		"Subject: " + os.Getenv("EMAIL_SUBJECT") + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(string(fmt.Sprintf("%+v", ec.ErrorObj))))

	if len(strings.TrimSpace(ec.Username)) != 0 {
		stmpAuth := smtp.PlainAuth("", ec.Username, ec.Password, ec.Host)

		err = smtp.SendMail(
			ec.Host+":"+ec.Port,
			stmpAuth,
			ec.Sender,
			ec.Recievers,
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
	// format reciever email
	for i := range ec.Recievers {
		ec.Recievers[i] = r.Replace(ec.Recievers[i])
		if err = conn.Rcpt(ec.Recievers[i]); err != nil {
			return err
		}
	}

	w, err := conn.Data()
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return conn.Quit()

}
