package views

import (
	"bytes"
)

type IView interface {
	GetTemplateName() string
	GetTemplateData() interface{}
	GetHTML() (*bytes.Buffer, error)
	GetText() (*bytes.Buffer, error)
}

type ViewData struct {
	Data interface{}
}

type BaseView struct {
	Template string
	Data     interface{}
}

func (v *BaseView) GetTemplateName() string {
	return v.Template
}

func (v *BaseView) GetTemplateData() interface{} {
	return v.Data
}

func (v *BaseView) GetHTML() (*bytes.Buffer, error) {
	return GetHTMLView(v)
}

func (v *BaseView) GetText() (*bytes.Buffer, error) {
	return GetTextView(v)
}
