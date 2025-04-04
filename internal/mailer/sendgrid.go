package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey string, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{fromEmail: fromEmail, apiKey: apiKey, client: client}
}

func (s *SendGridMailer) Send(templateFile string, username string, email string, data any, isSanbox bool) (int, error) {
	from := mail.NewEmail(FromName, s.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	htmlContent := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlContent, "body", data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", htmlContent.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSanbox,
		},
	})

	for i := range maxRetries {
		response, err := s.client.Send(message)
		if err != nil {
			// exponential backoff
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		return response.StatusCode, nil
	}

	return -1, fmt.Errorf("failed to send email after %d retries", maxRetries)
}
