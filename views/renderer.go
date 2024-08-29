package views

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Renderer struct {
	env      string
	rootPath string
	tmplPath string
	files    embed.FS
}

var envRender = "development"

var Manifest = map[string]string{}

var TemplateCache = map[string]*template.Template{}

func New(env string, rootPath string, files embed.FS) *Renderer {
	renderer := Renderer{
		env:      env,
		rootPath: rootPath,
		tmplPath: "app/templates",
		files:    files,
	}

	envRender = env

	return &renderer
}

func (ren *Renderer) Render(w http.ResponseWriter, r *http.Request, tmpl string, data any) {
	parsedTmpl, err := ren.fetchTemplate(tmpl)

	if err != nil {
		log.Println(err)
		return
	}

	err = parsedTmpl.Execute(w, nil)

	if err != nil {
		log.Println("Error on template exec:", err)
		return
	}

	// log.Println("rendering ...")
}

func (ren *Renderer) fetchTemplate(tmpl string) (*template.Template, error) {
	if ren.env != "development" {
		parsedTmpl, err := ren.templateFromCache(tmpl)

		return parsedTmpl, err
	}

	parsedTmpl, err := ren.templateFromDisk(tmpl)

	return parsedTmpl, err
}

func (ren *Renderer) templateFromCache(tmpl string) (*template.Template, error) {
	ren.loadManifest()

	parsedTmpl, ok := TemplateCache[tmpl]

	if ok {
		return parsedTmpl, nil
	}

	log.Println("Template cache is empty, loading cache")

	err := ren.templateCache()

	if err != nil {
		return nil, err
	}

	parsedTmpl, ok = TemplateCache[tmpl]

	if ok {
		return parsedTmpl, nil
	}

	return nil, errors.New(fmt.Sprintf("Template '%s' not found", tmpl))
}

func (ren *Renderer) templateCache() error {
	var err error

	pages, shared := ren.fetchTemplatefiles()

	fmt.Println("Main templates:")
	fmt.Println(strings.Join(pages, "\n"))
	fmt.Println("Base templates:")
	fmt.Println(strings.Join(shared, "\n"))

	for _, page := range pages {
		name := filepath.Base(page)
		key := ren.templateKey(page)

		// tmpl, err := template.New(name).ParseFiles(page)
		tmpl := template.Must(template.New(name).Funcs(funcMap).ParseFS(ren.files, page))

		if len(shared) > 0 {
			tmpl, err = tmpl.ParseFS(ren.files, shared...)
			if err != nil {
				return err
			}
		}

		// fmt.Println("name:", key)
		// fmt.Println("tmpl:", tmpl)
		TemplateCache[key] = tmpl
	}

	return nil
}

func (ren *Renderer) fetchTemplatefiles() (pages []string, shared []string) {
	pages = []string{}
	shared = []string{}

	filePattern := ren.tmplPath + "/*/*.tmpl"

	// get all the template files
	files, err := fs.Glob(ren.files, filePattern)

	if err != nil {
		log.Println("Error finding tempates:", err)
		return pages, shared
	}

	for _, file := range files {
		if strings.Index(file, "/layouts/") >= 0 {
			shared = append(shared, file)
		} else {
			pages = append(pages, file)
		}
	}

	return pages, shared
}

func (ren *Renderer) templateKey(file string) string {
	prefix := ren.tmplPath + "/"

	return strings.TrimPrefix(file, prefix)
}

func (ren *Renderer) templateFromDisk(tmpl string) (*template.Template, error) {
	ren.loadManifestDev()

	tmplParts := []string{ren.rootPath, ren.tmplPath, tmpl}
	sharedParts := []string{ren.rootPath, ren.tmplPath, "/layouts/*.tmpl"}

	tmplFile := strings.Join(tmplParts, "/")
	sharedPattern := strings.Join(sharedParts, "/")

	name := filepath.Base(tmplFile)

	// tmpl, err := template.New(name).ParseFiles(page)
	parsedTmpl := template.Must(template.New(name).Funcs(funcMap).ParseFiles(tmplFile))

	// get all shared template files
	shared, err := filepath.Glob(sharedPattern)

	if err != nil {
		log.Println("Error finding tempates:", err)
		return nil, err
	}

	if len(shared) > 0 {
		parsedTmpl, err = parsedTmpl.ParseFiles(shared...)
		if err != nil {
			return nil, err
		}
	}

	log.Println("Templates:")
	fmt.Println(strings.Join(shared, "\n"))
	fmt.Println(tmplFile)

	return parsedTmpl, nil
}

func (ren *Renderer) loadManifest() error {
	mfile, err := ren.files.ReadFile("public/assets/manifest.json")

	if err != nil {
		log.Println("Error reading assets manifest.json:", err)
	}

	err = json.Unmarshal([]byte(mfile), &Manifest)

	if err != nil {
		log.Println("Error parsing assets manifest.json:", err)
	}

	return nil
}

func (ren *Renderer) loadManifestDev() error {
	mpath := ren.rootPath + "/frontend/builds/manifest.json"

	mfile, err := os.ReadFile(mpath)

	if err != nil {
		log.Println("Error reading assets manifest.json:", err)
	}

	err = json.Unmarshal([]byte(mfile), &Manifest)

	if err != nil {
		log.Println("Error parsing assets manifest.json:", err)
	}

	return nil
}
