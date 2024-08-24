package templating

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Renderer struct {
	env string
	rootPath string
	viewPath string
	files    embed.FS
}


func New (env string, rootPath string, files embed.FS) *Renderer {
	renderer := Renderer{
		env: env,
		rootPath: rootPath,
		viewPath: "app/views",
		files: files,
	}

	return &renderer
}

func (ren *Renderer) Render (w http.ResponseWriter, r *http.Request, view string, data any) {

	tmpl, err := ren.parseTemplate(view)

	if err != nil {
		log.Println("Error on:", err)
		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		log.Println("Error on:", err)
		return
	}

	log.Println("rendering ...")
}

func (ren *Renderer) parseTemplate(view string) (*template.Template, error) {
	viewParts   := []string{ren.viewPath, view}
	layoutParts := []string{ren.viewPath, "layouts/application.tmpl"}

	viewFile   := strings.Join(viewParts, "/")
	layoutFile := strings.Join(layoutParts, "/")

	if ren.env != "development" {
		tmpl, err := template.ParseFS(ren.files, viewFile, layoutFile)

		return tmpl, err
	}

	tmpl, err := ren.parseTemplateDev(view)

	return tmpl, err
}


func (ren *Renderer) parseTemplateDev(view string) (*template.Template, error) {
	viewParts   := []string{ren.rootPath, ren.viewPath, view}
	layoutParts := []string{ren.rootPath, ren.viewPath, "layouts/application.tmpl"}

	viewFile   := strings.Join(viewParts, "/")
	layoutFile := strings.Join(layoutParts, "/")

	tmpl, err := template.ParseFiles(viewFile, layoutFile)

	return tmpl, err
}
