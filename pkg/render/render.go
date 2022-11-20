package render

import (
	"bytes"
	"fmt"
	"hello-world-server/pkg/config"
	"hello-world-server/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func addDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// NewTemplates sets the config for the templates package.
func NewTemplates(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders templates using html/template.
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTempleteCache()
	}

	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("Tempalte Not Found in the cache try to refresh the cache and try again.")
	}

	buf := new(bytes.Buffer)
	td = addDefaultData(td)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing the template to the browser", err)
	}

}

// CreateTempleteCache create a template map as a cache
func CreateTempleteCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.*")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		matches, err := filepath.Glob("./templates/*.layout.*")
		if err != nil {
			return nil, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.*")
			if err != nil {
				return nil, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
