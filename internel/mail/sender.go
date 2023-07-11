package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name          string
	fromEmailAddr string
	fromEmailPass string
}

func NewGmailSender(
	name string,
	fromEmailAddr string,
	fromEmailPassword string,
) EmailSender {
	return &GmailSender{
		name:          name,
		fromEmailAddr: fromEmailAddr,
		fromEmailPass: fromEmailPassword,
	}
}

func (mailSender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>",
		mailSender.name, mailSender.fromEmailAddr)
	e.Subject = subject
	e.HTML = []byte(content)
	e.Cc = cc
	e.To = to
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}
	smtpAuth := smtp.PlainAuth("", mailSender.fromEmailAddr, mailSender.fromEmailPass, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
