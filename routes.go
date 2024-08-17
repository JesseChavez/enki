package enki

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (c *Enki) defaultRoutes() *Mux {
	mux := chi.NewRouter()

	// adding to log request details
	mux.Use(middleware.Logger)

	// adding request ID
	mux.Use(middleware.RequestID)

	// adding to recover panics gracefully and log stack trace
	mux.Use(middleware.Recoverer)

	// adding server status check
	mux.Use(middleware.Heartbeat("/ping"))

	return mux
}
