package mailer

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunService struct {
	DefaultSender     string
	DefaultSenderName string
	URL               string
	PrivateKey        string
}

func (mgs *MailgunService) Send(ctx context.Context, mail Mail) (*MailResponse, error) {
	mg := mailgun.NewMailgun(mgs.URL, mgs.PrivateKey)

	sender := fmt.Sprintf("%s <%s>", mail.SenderName, mail.Sender)
	recipients := []string{}
	for _, r := range mail.Recipients {
		recipients = append(recipients, fmt.Sprintf("%s <%s>", r.Name, r.Email))
	}
	message := mg.NewMessage(sender, mail.Subject, mail.BodyText, recipients...)
	message.SetHtml(mail.BodyHTML)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)
	if err != nil {
		return nil, err
	}

	return &MailResponse{
		ID:       id,
		Response: resp,
	}, nil
}
