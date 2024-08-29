package enki

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (ek *Enki) InitRouting() *Mux {
	trunkMux, contextMux := ek.appRoutes()

	ek.Routes = trunkMux

	return contextMux
}

func (ek *Enki) NewRouter() *Mux {
	mux := chi.NewRouter()

	return mux
}

func (ek *Enki) appRoutes() (*Mux, *Mux) {
	trunkMux := chi.NewRouter()

	// adding to log request details
	trunkMux.Use(middleware.Logger)

	// adding request ID
	trunkMux.Use(middleware.RequestID)

	// adding to recover panics gracefully and log stack trace
	trunkMux.Use(middleware.Recoverer)

	// adding server status check
	trunkMux.Use(middleware.Heartbeat("/ping"))

	// NOTE: chi does not like mount static file server in sub route
	ek.staticAssets(contextPath, trunkMux)

	if contextPath == "/" {
		return trunkMux, trunkMux
	}

	valid, err := regexp.MatchString(`^/[a-zA-Z]{1}[a-zA-Z-_0-9]+$`, contextPath)

	if err != nil {
		log.Fatal(err.Error())
	}

	if !valid {
		log.Fatalf("Invalid ContextPath %v, it need to be like \"/api\"", contextPath)
	}

	contextMux := chi.NewRouter()

	trunkMux.Mount(contextPath, contextMux)

	return trunkMux, contextMux
}

func (ek *Enki) staticAssets(contextPath string, mux *Mux) {
	prefixPath := ""

	if contextPath != "/" {
		prefixPath = contextPath
	}

	indexPath := fmt.Sprintf("%s/assets/", prefixPath)
	stripPath := fmt.Sprintf("%s/assets", prefixPath)
	handlePath := fmt.Sprintf("%s/assets/*", prefixPath)



	// Disable assets index page
	mux.Get(indexPath, func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusNotFound)
	})

	if ek.Env != "development" {
		// Handle static assets for production or test env.
		mux.Handle(handlePath, http.StripPrefix(stripPath, ek.staticHandler()))
	} else {
		// Handle static assets for development env.
		mux.Handle(handlePath, http.StripPrefix(stripPath, ek.staticHandlerDev()))
		// mux.Handle("/assets/*", http.StripPrefix("/assets", ek.staticHandlerDev()))
	}
}

func (ek *Enki) staticHandler() http.HandlerFunc {
	publicDir, err := fs.Sub(Resources, "public/assets")

	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(publicDir))

	return func(w http.ResponseWriter, r *http.Request) {
		// add custom headers
		w.Header().Add("Cache-Control", "public, max-age=31536000, immutable")
		// w.Header().Add("X-Frame-Options", "SAMEORIGIN")
		// w.Header().Add("X-Content-Type-Options", "nosniff")
		// w.Header().Add("X-XSS-Protection", "1; mode=block")
		// w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubdomains;")
		// w.Header().Add("Referrer-Policy", "no-referrer-when-downgrade")
		w.Header().Add("X-API-Assets", "static")

		fileServer.ServeHTTP(w, r)
	}
}

func (ek *Enki) staticHandlerDev() http.HandlerFunc {
	publicDir := rootPath + "/frontend/builds"

	if  _, err := os.Stat(publicDir); os.IsNotExist(err) {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir(publicDir))

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-API-Assets", "static")

		fileServer.ServeHTTP(w, r)
	}
}
