package mail

import (
	"crypto/tls"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	d          *gomail.Dialer
	Sender     string
	APIAddress string
	Connected  bool
}

func New(addr, login, password, apiaddr string, port int) *Mailer {
	d := gomail.NewDialer(addr, port, login, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec

	mes := gomail.NewMessage()

	mes.SetHeader("From", login)
	mes.SetHeader("To", login)
	mes.SetHeader("Subject", "New instance entry initialized!")
	mes.SetBody("text/html", "Hello from exchanger, new instance entry have been launched with following email!")

	err := d.DialAndSend(mes)
	if err != nil {
		logrus.Errorf("unable to send email notification: %v", err)
	}
	return &Mailer{
		d:          d,
		Sender:     login,
		APIAddress: apiaddr,
		Connected:  err == nil,
	}
}
