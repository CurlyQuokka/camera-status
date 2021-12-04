package mailer

import (
	"fmt"
	"net/smtp"
)

const (
	ErrSubject = "Camera ERROR!"
	OkSubject  = "Camera is fine!"
)

type Mailer struct {
	sender   string
	receiver string
	secToken string
	smtpHost string
	smtpPort string
}

func NewMailer(from, to, securityToken, smtpSrv, smtpSrvPort string) *Mailer {
	return &Mailer{
		sender:   from,
		receiver: to,
		secToken: securityToken,
		smtpHost: smtpSrv,
		smtpPort: smtpSrvPort,
	}
}

func (m *Mailer) SendMail(subject, message string) error {
	msg := []byte("To: " + m.sender + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")

	auth := smtp.PlainAuth("", m.sender, m.secToken, m.smtpHost)

	toSlice := []string{m.receiver}

	err := smtp.SendMail(m.smtpHost+":"+m.smtpPort, auth, m.sender, toSlice, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
