package app

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pages, err := filepath.Glob(path + "/ui/templates/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		t, err := template.New(name).Funcs(functions).ParseFiles(path + "/ui/templates/base.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		t, err = t.ParseGlob(path + "/ui/templates/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the  page template.
		t, err = t.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map as normal...
		cache[name] = t
	}

	return cache, nil
}

func NewTemplateCacheForTests() (map[string]*template.Template, error) {
	// fmt.Println(os.Getwd())
	cache := map[string]*template.Template{}

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path = strings.TrimSuffix(path, "/internal/handlers")

	pages, err := filepath.Glob(path + "/ui/templates/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		t, err := template.New(name).Funcs(functions).ParseFiles(path + "/ui/templates/base.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		t, err = t.ParseGlob(path + "/ui/templates/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the  page template.
		t, err = t.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map as normal...
		cache[name] = t
	}

	return cache, nil
}
