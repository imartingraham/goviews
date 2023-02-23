package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strings"
	"views/pkg/mailer"
	"views/pkg/views"
)

//go:embed templates/*
var content embed.FS

type User struct {
	ID        int
	FirstName string
	LastName  string
}

type Button struct {
	Title string
	Type  string
	HRef  string
}

type UserEmail struct {
	User   User
	Button Button
}

func main() {
	err := mailer.Configure(&mailer.Mailer{
		DefaultSenderName: "Ian Graham",
		DefaultSender:     "ian@iangraham.io",
		TemplateConfig: &views.ViewConfig{
			Directory:            "templates",
			DefinitionsDirectory: "definitions",
			Content:              content,
			FuncMap: map[string]interface{}{
				"title": strings.Title,
			},
		},
		// Service: &mailer.SendgridService{
		// 	APIKey: os.Getenv("SENDGRID_API_KEY"),
		// },
		Service: &mailer.MailgunService{
			URL:        os.Getenv("MAILGUN_DOMAIN"),
			PrivateKey: os.Getenv("MAILGUN_PRIVATE_KEY"),
		},
	})

	if err != nil {
		panic(err)
	}

	mail := mailer.Mail{
		Recipients: []mailer.Recipient{
			{
				Name:  "Testing Test",
				Email: "ian+testing@iangraham.io",
			},
		},
		Subject:  "Testing",
		Template: "email/index",
		Data: UserEmail{
			User: User{
				ID:        1,
				FirstName: "test",
				LastName:  "Testerson",
			},
			Button: Button{
				HRef:  "https://investorkeep.com",
				Title: "Hello",
				Type:  "submit",
			},
		},
	}
	ctx := context.Background()
	resp, err := mailer.Send(ctx, mail)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nResp: %v\n", resp)

	err = views.Configure(&views.ViewConfig{
		Directory:            "templates",
		DefinitionsDirectory: "definitions",
		Content:              content,
		FuncMap: map[string]interface{}{
			"title": strings.Title,
		},
	})

	if err != nil {
		panic(err)
	}

	user := User{
		ID:        1,
		FirstName: "test",
		LastName:  "Testerson",
	}
	testView := views.BaseView{
		Template: "email/index",
		Data: UserEmail{
			User: user,
			Button: Button{
				Title: "Hello",
				Type:  "submit",
			},
		},
	}

	html, err := testView.GetHTML()
	if err != nil {
		panic(err)
	}
	txt, err := testView.GetText()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nTxt: %v\n", txt)
	fmt.Printf("\nHTML: %v\n", html)
}
