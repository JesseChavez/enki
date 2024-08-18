package enki

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
)

const version = "0.0.1"

type Mux = chi.Mux

var ContextPath = "/"

type Enki struct {
	AppName string
	trunk   *Mux
	Routes  *Mux
}

func New(name string) Enki {
	app := Enki{}

	app.AppName = name

	app.trunk, app.Routes = app.appRoutes()

	return app
}

func (enki *Enki) ListenAndServe(port string) {
	webPort := fmt.Sprintf(":%v", port)

	server := &http.Server{
		Addr:         webPort,
		Handler:      enki.trunk,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	log.Println("Web Applications is starting...")
	log.Println("* Enki version:", enki.Version())
	log.Println("*   Go version:", runtime.Version())
	log.Println("*   Process ID:", os.Getpid())
	log.Println("*   Using port:", port)
	log.Println("Use Ctrl-C to stop")

	log.Fatal(server.ListenAndServe())
}

func (enki *Enki) Version() string {
	return version
}
