package templatemanager

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
)

type TemplatesPath struct {
	Pages string
	Partials string
	Layout string
}

type TemplateManager struct {
	Functions template.FuncMap
	TemplatesPath TemplatesPath
	Fsys fs.FS
}

func NewTemplateManager(tp TemplatesPath, fsys fs.FS) *TemplateManager {
	return &TemplateManager{
		TemplatesPath: tp,
		Functions: template.FuncMap{},
		Fsys: fsys,
	}
}

var templates map[string]*template.Template

func (t *TemplateManager) GetTemplates() (map[string]*template.Template, error) {
	if templates == nil {
		var err error
		if t.Fsys != nil {
			templates, err = t.ParseTemplatesFS(t.Fsys)
			if err != nil {
				return nil, err
			}

			return templates, nil
		}


		templates, err = t.ParseTemplates()
		if err != nil {
			return nil, err
		}
	}

	return templates, nil
}

func (t *TemplateManager) ParseTemplates() (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}
	pages, err := filepath.Glob(t.TemplatesPath.Pages)
	if err != nil {
		return nil, err
	}

	partials, err := filepath.Glob(t.TemplatesPath.Partials)
	if err != nil {
		return nil, err
	}

	for _, p := range pages {
		name := filepath.Base(p)
		compound := append([]string{t.TemplatesPath.Layout, p}, partials...)
		ts, err := template.New(name).Funcs(t.Functions).ParseFiles(compound...)

		if err != nil {
			return nil, err
		}

		templates[name] = ts
	}

	return templates, nil
}

func (t *TemplateManager) ParseTemplatesFS(fsys fs.FS) (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}
	pages, err := filepath.Glob(t.TemplatesPath.Pages)
	if err != nil {
		return nil, err
	}

	for _, p := range pages {
		name := filepath.Base(p)
		ts, err := template.New(name).Funcs(t.Functions).ParseFS(fsys, []string{
			t.TemplatesPath.Partials,
			t.TemplatesPath.Layout,
			p,
		}...)

		if err != nil {
			return nil, err
		}

		templates[name] = ts
	}

	return templates, nil
}

func (t *TemplateManager) Generate(w io.Writer, name string, data interface{}) error {
	templates, err := t.GetTemplates()
	if err != nil {
		return err
	}

	ts, ok := templates[name]
	if !ok {
		return fmt.Errorf("template does not exist: %s", name)
	}

	// execute in buffer to catch errors early
	buf := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buf, "layout", data); err != nil {
		return err
	}

	_, err = buf.WriteTo(w)

	return err
}