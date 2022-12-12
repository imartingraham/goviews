package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

// ViewManager handles loading and parsing view templates
type ViewManager struct {
	content              fs.FS
	Directory            string
	DefinitionsDirectory string
	Templates            map[string]*template.Template
	Definitions          map[string][]string
}

// ViewConfig is used to configure how ikviews loads templates
type ViewConfig struct {
	Directory            string
	DefinitionsDirectory string
	Content              fs.FS
}

// ikViewManager is the ViewManager singleton
var ikViewManager = ViewManager{
	Directory:            "templates",
	DefinitionsDirectory: "definitions",
	Templates:            map[string]*template.Template{},
	Definitions:          map[string][]string{},
}

// Configure sets the Directory, DefinitionsDirectory and Content
// values and then loads all the templates found in `Content`
func Configure(vc *ViewConfig) error {
	if vc.Content == nil {
		return fmt.Errorf("ViewConfig.Content must contain an instance of `embed.FS` ")
	}
	ikViewManager.content = vc.Content
	if vc.Directory != "" {
		ikViewManager.Directory = vc.Directory
	}

	if vc.DefinitionsDirectory != "" {
		ikViewManager.DefinitionsDirectory = vc.DefinitionsDirectory
	}

	return ikViewManager.loadTemplates()
}

// GetHTMLView returns a populated that matches the name with an .html extension
func GetHTMLView(name string, data interface{}) (*bytes.Buffer, error) {
	return ikViewManager.getPopulatedTemplate(name, "html", data)
}

// GetTextView returns a populated that matches the name with an .txt extension
func GetTextView(name string, data interface{}) (*bytes.Buffer, error) {
	return ikViewManager.getPopulatedTemplate(name, "txt", data)
}

// getTemplate returns the cached template
func (vm *ViewManager) getTemplate(name string, format string) (*template.Template, error) {
	tmplName := vm.Directory + "/" + name + "." + format
	tmpl, ok := vm.Templates[tmplName]
	if !ok {
		return nil, fmt.Errorf("template file at %s not found", tmplName)
	}
	return tmpl, nil
}

// getPopulatedTemplate injects the data into the template
func (vm *ViewManager) getPopulatedTemplate(name string, format string, data interface{}) (*bytes.Buffer, error) {
	tmpl, err := vm.getTemplate(name, format)
	if err != nil {
		return nil, err
	}

	tpl := &bytes.Buffer{}
	err = tmpl.Execute(tpl, data)
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

// loadDefinitions loads template definitions to be used when creating templates
func (vm *ViewManager) loadDefinitions() error {
	return fs.WalkDir(vm.content, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.Contains(path, vm.DefinitionsDirectory) {
			ext := filepath.Ext(path)
			_, ok := vm.Definitions[ext]
			if !ok {
				vm.Definitions[ext] = []string{}
			}
			vm.Definitions[ext] = append(vm.Definitions[ext], path)
		}
		return nil
	})
}

// loadTemplates loads all of the templates and definition files into a template.
// Loading the templates upfront allows us to catch any syntax or loading errors
// rather than catching errors at runtime
func (vm *ViewManager) loadTemplates() error {
	err := vm.loadDefinitions()
	if err != nil {
		return err
	}

	return fs.WalkDir(vm.content, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || strings.Contains(path, vm.DefinitionsDirectory) {
			return nil
		}

		if err != nil {
			return err
		}
		ext := filepath.Ext(path)

		// only load the definitions for the extention
		paths := append(vm.Definitions[ext], path)

		var tmpl *template.Template
		tmpl, err = template.New(filepath.Base(d.Name())).ParseFS(vm.content, paths...)
		if err != nil {
			return err
		}
		vm.Templates[path] = tmpl

		return nil
	})
}
