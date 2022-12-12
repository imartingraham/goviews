package main

import (
	"embed"
	"fmt"
	"views/pkg/views"
)

//go:embed templates/*
var content embed.FS

type User struct {
	ID        int
	FirstName string
	LastName  string
}

func main() {
	err := views.Configure(&views.ViewConfig{
		Directory:            "templates",
		DefinitionsDirectory: "definitions",
		Content:              content,
	})

	if err != nil {
		panic(err)
	}

	data := User{
		ID:        1,
		FirstName: "Test",
		LastName:  "Testerson",
	}
	testView := views.BaseView{
		Template: "email/index",
		Data:     data,
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
