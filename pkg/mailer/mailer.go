package mailer

import (
	"net/smtp"

	"github.com/CurlyQuokka/camera-status/pkg/logger"
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
	logger   *logger.Logger
}

func NewMailer(from, to, securityToken, smtpSrv, smtpSrvPort string, log logger.Logger) *Mailer {
	return &Mailer{
		sender:   from,
		receiver: to,
		secToken: securityToken,
		smtpHost: smtpSrv,
		smtpPort: smtpSrvPort,
		logger:   &log,
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
		(*m.logger).Error(err.Error())
		return err
	}
	infoMsg := "Sent message: " + message
	(*m.logger).Info(infoMsg)
	return nil
}
