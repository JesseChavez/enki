package renderer

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/JesseChavez/enki/logger"
)

type Renderer struct {
	env        string
	prefixPath string
	rootPath   string
	tmplPath   string
	files      embed.FS
	Log        logger.ILogger
}

var prefixPath = ""

var xmlDeclaration = `<?xml version="1.0" encoding="utf-8"?>` + "\n"

var envRender string

var Manifest = map[string]string{}

var TemplateCache = map[string]*template.Template{}

func New(env string, logger logger.ILogger, contextPath string, rootPath string, files embed.FS) *Renderer {
	if contextPath != "/" {
		prefixPath = contextPath
	}

	renderer := Renderer{
		env:        env,
		prefixPath: prefixPath,
		rootPath:   rootPath,
		tmplPath:   "app/templates",
		files:      files,
		Log:        logger,
	}

	envRender = env

	return &renderer
}

func (ren *Renderer) Render(w http.ResponseWriter, status int, meta map[string]string, data any) {
	tmpl := meta["template"]

	parsedTmpl, err := ren.fetchTemplate(tmpl)

	if err != nil {
		ren.Log.Error("Render, fetch template:", err)
		return
	}

	w.Header().Set("Content-Type", ren.contentTypeHeader(meta["mime_type"], meta["charset"]))

	w.WriteHeader(status)

	err = parsedTmpl.Execute(w, data)

	if err != nil {
		ren.Log.Error("Render, template execute", err)
		return
	}

	// log.Println("rendering ...")
}

 func (ren *Renderer)RenderJSON(w http.ResponseWriter, status int, data any) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	err := enc.Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	w.Write(buf.Bytes())
}

func (ren *Renderer) RenderHTML(w http.ResponseWriter, status int, tmpl string, data any) {
	parsedTmpl, err := ren.fetchTemplate(tmpl)

	if err != nil {
		ren.Log.Error("RenderHTML, fetch template:", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(status)

	err = parsedTmpl.Execute(w, data)

	if err != nil {
		ren.Log.Error("RenderHTML, template execute", err)
		return
	}

	// log.Println("rendering ...")
}

func (ren *Renderer) RenderXML(w http.ResponseWriter, status int, tmpl string, data any) {
	parsedTmpl, err := ren.fetchTemplate(tmpl)

	if err != nil {
		ren.Log.Error("RenderXML, fetch template:", err)
		return
	}

	buf := new(bytes.Buffer)

	err = parsedTmpl.Execute(buf, data)

	if err != nil {
		ren.Log.Error("RenderXML, template execute", err)
		return
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	w.WriteHeader(status)

	w.Write([]byte(xmlDeclaration + buf.String()))

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

	ren.Log.Info("Template cache is empty, loading cache")

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

	ren.Log.Debug("Main templates:")
	ren.Log.Debug(strings.Join(pages, "\n"))
	ren.Log.Debug("Base templates:")
	ren.Log.Debug(strings.Join(shared, "\n"))

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
		ren.Log.Error("Error finding tempates:", err)
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
		ren.Log.Error("Error finding tempates:", err)
		return nil, err
	}

	if len(shared) > 0 {
		parsedTmpl, err = parsedTmpl.ParseFiles(shared...)
		if err != nil {
			return nil, err
		}
	}

	ren.Log.Debug("Templates:")
	ren.Log.Debug(strings.Join(shared, "\n"))
	ren.Log.Debug(tmplFile)

	return parsedTmpl, nil
}

func (ren *Renderer) loadManifest() error {
	mfile, err := ren.files.ReadFile("public/assets/manifest.json")

	if err != nil {
		ren.Log.Error("Error reading assets manifest.json:", err)
	}

	err = json.Unmarshal([]byte(mfile), &Manifest)

	if err != nil {
		ren.Log.Error("Error parsing assets manifest.json:", err)
	}

	return nil
}

func (ren *Renderer) loadManifestDev() error {
	mpath := ren.rootPath + "/tmp/assets/manifest.json"

	mfile, err := os.ReadFile(mpath)

	if err != nil {
		ren.Log.Error("Error reading assets manifest.json:", err)
	}

	err = json.Unmarshal([]byte(mfile), &Manifest)

	if err != nil {
		ren.Log.Error("Error parsing assets manifest.json:", err)
	}

	return nil
}

func (ren *Renderer) contentTypeHeader(mimeType string, charset string) string {
	contentType := mimeType + "; charset=" + charset

	return contentType
}
