package view

import (
	"embed"
	"fmt"
	"net/http"
	"github.com/JesseChavez/enki/renderer"
	"github.com/go-chi/chi/v5"
)

type ViewSupport struct {
	env        string
	prefixPath string
	Renderer   *renderer.Renderer
}

var defaultMimeType = "text/html"

var defaultCharset = "utf-8"

func New(env string, contextPath string, rootPath string, files embed.FS) *ViewSupport {
	prefixPath := ""

	if contextPath != "/" {
		prefixPath = contextPath
	}

	support := ViewSupport{
		env:        env,
		prefixPath: prefixPath,
		Renderer: renderer.New(env, contextPath, rootPath, files),
	}

	return &support
}

func (supp *ViewSupport) Render(w http.ResponseWriter, status int, view *ActionView) {
	meta := make(map[string]string)

	meta["template"] = view.Template
	
	meta["mime_type"] = supp.defaultValue(view.MimeType, defaultMimeType)
	meta["charset"]   = supp.defaultValue(view.Charset, defaultCharset)

	supp.Renderer.Render(w, status, meta, view.Data)
}

func (supp *ViewSupport) RenderHTML(w http.ResponseWriter, status int, view *ActionView) {
	template := view.Template

	supp.Renderer.RenderHTML(w, status, template, view.Data)
}

func (supp *ViewSupport) RenderXML(w http.ResponseWriter, status int, view *ActionView) {
	template := view.Template

	supp.Renderer.RenderXML(w, status, template, view.Data)
}

func (supp *ViewSupport) RoutePath(path string) string {
	return fmt.Sprintf("%s/%v", supp.prefixPath, path)
}

func (supp *ViewSupport) AssetPath(fileName string) string {
	return fmt.Sprintf("%s/assets/%v", supp.prefixPath, fileName)
}

// URLParam returns the url parameter from a http.Request object.
func (supp *ViewSupport) URLParam(r *http.Request, key string) string {
	value := chi.URLParam(r, key)

	return value
}

func (supp *ViewSupport) defaultValue(value string, fallback string) string {
	if value == "" {
		value = fallback
	}

	return value
}
