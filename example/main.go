package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jsainthilaire/templatemanager"
)

//go:embed "templates"
var embeddedTemplates embed.FS

func main() {
	// example using absolute paths
	templatesPath, err := filepath.Abs("example")
	if err != nil {
		// handle err
	}

	tp := templatemanager.TemplatesPath{
		Layout: fmt.Sprintf("%s/templates/layout.html", templatesPath),
		Pages: fmt.Sprintf("%s/templates/pages/*.html", templatesPath),
		Partials: fmt.Sprintf("%s/templates/partials/*.html", templatesPath),
	}

	t := templatemanager.NewTemplateManager(tp, nil)
	err = t.Generate(os.Stdout, "page1.html", map[string]string{
		"Page1Data": "page 1 data",
		"PartialData": "partial/piece data",
	})
	if err != nil {
		// handle error
	}

	err = t.Generate(os.Stdout, "page2.html", map[string]string{
		"Page2Data": "page 2 data",
		"PartialData": "partial/piece data",
	})

	if err != nil {
		// handle error
	}

	// example using embedded files
	tpFS := templatemanager.TemplatesPath{
		Layout: "templates/layout.html",
		Pages: "templates/pages/*.html",
		Partials: "templates/partials/*.html",
	}

	tFS := templatemanager.NewTemplateManager(tpFS, embeddedTemplates)
	err = tFS.Generate(os.Stdout, "page1.html", map[string]string{
		"Page1Data": "page 1 data",
		"PartialData": "partial/piece data",
	})
	if err != nil {
		// handle error
	}
}
