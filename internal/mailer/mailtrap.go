package mailer

import (
	"bytes"
	"fmt"
	"text/template"

	"gopkg.in/gomail.v2"
)

type MailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailtrapClient(fromEmail, apiKey string) (*MailtrapClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is missing")
	}
	return &MailtrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (mailer *MailtrapClient) Send(templateFile,
	username,
	email string,
	data any,
	isSandBox bool,
) error {

	template, err := template.ParseFS(FS, fmt.Sprintf("templates/%v", UserWelcomeTemplate))
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := template.ExecuteTemplate(subject, "subject", data); err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := template.ExecuteTemplate(body, "body", data); err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", mailer.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", mailer.apiKey)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
