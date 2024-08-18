package enki

import (
	"log"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (c *Enki) appRoutes() (*Mux, *Mux) {
	trunkMux := chi.NewRouter()

	// adding to log request details
	trunkMux.Use(middleware.Logger)

	// adding request ID
	trunkMux.Use(middleware.RequestID)

	// adding to recover panics gracefully and log stack trace
	trunkMux.Use(middleware.Recoverer)

	// adding server status check
	trunkMux.Use(middleware.Heartbeat("/ping"))

	if ContextPath == "/" {
		return trunkMux, trunkMux
	}

	valid, err := regexp.MatchString(`^/[a-zA-Z]{1}[a-zA-Z-_0-9]+$`, ContextPath)

	if err != nil {
		log.Fatal(err.Error())
	}

	if !valid {
		log.Fatalf("Invalid ContextPath %v, it need to be like \"/api\"", ContextPath)
	}

	contextMux := chi.NewRouter()

	trunkMux.Mount(ContextPath, contextMux)

	return trunkMux, contextMux
}
