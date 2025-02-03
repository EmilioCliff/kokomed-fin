package pkg

import (
	"bytes"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendMail(subject, content string, mimeType string, to, cc, bcc []string, attachmentNames []string, attachmentData [][]byte) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender GmailSender) SendMail(subject, content string, mimeType string, to, cc, bcc []string, attachmentNames []string, attachmentData [][]byte) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Bcc = bcc
	cc = append(cc, "clifftest33@gmail.com")
	e.Cc = cc

	for idx, attachFile := range attachmentData {
		attachmentReader := bytes.NewReader(attachFile)
		e.Attach(attachmentReader, attachmentNames[idx], mimeType)
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}