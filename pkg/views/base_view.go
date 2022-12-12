package views

import (
	"bytes"
)

type IView interface {
	GetTemplateData() interface{}
	GetTemplateName() string
	GetHTML() (*bytes.Buffer, error)
	GetText() (*bytes.Buffer, error)
}

type BaseView struct {
	Template string
	Data     interface{}
}

func (v *BaseView) GetTemplateData() interface{} {
	return v.Data
}

func (v *BaseView) GetTemplateName() string {
	return v.Template
}

func (v *BaseView) GetHTML() (*bytes.Buffer, error) {
	return GetHTMLView(v.Template, v.Data)
}

func (v *BaseView) GetText() (*bytes.Buffer, error) {
	return GetTextView(v.Template, v.Data)
}
