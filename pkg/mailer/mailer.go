package mailer

import (
	"context"
	"views/pkg/views"
)

type IMailerService interface {
	Send(context.Context, Mail) (*MailResponse, error)
}

type Mailer struct {
	DefaultSender     string
	DefaultSenderName string
	Service           IMailerService
	TemplateConfig    *views.ViewConfig
	viewManager       *views.ViewManager
}

type Recipient struct {
	Name  string
	Email string
}
type Mail struct {
	Sender     string
	SenderName string
	Subject    string
	Recipients []Recipient
	Template   string
	BodyText   string
	BodyHTML   string
	Data       interface{}
}

type MailResponse struct {
	ID       string
	Response string
}

var globalMailer *Mailer

func Configure(mailer *Mailer) error {
	m, err := Init(mailer)
	if err != nil {
		return err
	}
	globalMailer = m
	return nil
}

func Init(mailer *Mailer) (*Mailer, error) {
	viewManager, err := views.Init(mailer.TemplateConfig)
	if err != nil {
		return nil, err
	}

	mailer.viewManager = viewManager
	return mailer, nil
}

func Send(ctx context.Context, mail Mail) (*MailResponse, error) {
	return globalMailer.Send(ctx, mail)
}

func (m *Mailer) Send(ctx context.Context, mail Mail) (*MailResponse, error) {

	mailView := &views.BaseView{
		Template: mail.Template,
		Data:     mail.Data,
	}

	if mail.Sender == "" {
		mail.Sender = m.DefaultSender
	}
	if mail.SenderName == "" {
		mail.SenderName = m.DefaultSenderName
	}

	html, err := m.viewManager.GetPopulatedTemplate(mailView, "html")
	if err != nil {
		return nil, err
	}
	mail.BodyHTML = html.String()

	txt, err := m.viewManager.GetPopulatedTemplate(mailView, "txt")
	if err != nil {
		return nil, err
	}
	mail.BodyText = txt.String()

	return m.Service.Send(ctx, mail)
}
