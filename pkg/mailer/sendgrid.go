package mailer

import (
	"context"

	"github.com/sendgrid/sendgrid-go"
	smail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridService struct {
	DefaultSender     string
	DefaultSenderName string
	APIKey            string
}

func (ss *SendgridService) Send(ctx context.Context, mail Mail) (*MailResponse, error) {

	from := smail.NewEmail(mail.SenderName, mail.Sender)
	m := smail.NewV3Mail()
	m.SetFrom(from)
	html := smail.NewContent("text/html", mail.BodyHTML)
	txt := smail.NewContent("text/plain", mail.BodyText)
	// For sendgrid the order of plain and text html is important.
	// text/plain must be first...
	m.AddContent(txt, html)

	personalization := smail.NewPersonalization()
	tos := []*smail.Email{}
	for _, r := range mail.Recipients {
		tos = append(tos, smail.NewEmail(r.Name, r.Email))
	}
	personalization.AddTos(tos...)
	personalization.Subject = mail.Subject
	m.AddPersonalizations(personalization)

	client := sendgrid.NewSendClient(ss.APIKey)
	resp, err := client.Send(m)
	if err != nil {
		return nil, err
	}

	return &MailResponse{
		Response: resp.Body,
	}, nil
}
