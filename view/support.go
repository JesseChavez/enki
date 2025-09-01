package view

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/JesseChavez/enki/renderer"
	"github.com/go-chi/chi/v5"
)

type ViewSupport struct {
	env        string
	api        bool
	csr        bool
	prefixPath string
	Renderer   *renderer.Renderer
}

type ViewSpec struct {
	SpecID string
	Name   string
	Data   template.JS
}

var defaultMimeType = "text/html"

var defaultCharset = "utf-8"

func New(env string, api bool, csr bool, contextPath string, rootPath string, files embed.FS) *ViewSupport {
	prefixPath := ""

	if contextPath != "/" {
		prefixPath = contextPath
	}

	support := ViewSupport{
		env:        env,
		api:        api,
		csr:        csr,
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

	if supp.csr {
		supp.Renderer.Render(w, status, meta, marshalData(view))
	} else {
		supp.Renderer.Render(w, status, meta, view.Data)
	}
}

func (supp *ViewSupport) RenderJSON(w http.ResponseWriter, status int, data any) {
	supp.Renderer.RenderJSON(w, status, data)
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

func marshalData(view *ActionView) any {
	spec := ViewSpec{
		SpecID: view.Name + "-" + randomHex(8),
		Name: view.Name,
	}

	out, err := json.MarshalIndent(view.Data, "", "  ")

	data := "{}"

	if err == nil {
		data = string(out)
	}

	if view.Debug {
		log.Printf("Data: %s\n", data)
	}

	spec.Data = template.JS(data)

	return &spec
}

func randomHex(n int) (string) {
	bytes := make([]byte, n)

	_, err := rand.Read(bytes);

	if err != nil {
		log.Println("Error generation random hex:", err)
		return "0"
	}

	return hex.EncodeToString(bytes)
}
