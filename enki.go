package enki

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/JesseChavez/enki/database"
	"github.com/JesseChavez/spt"
	"github.com/go-chi/chi/v5"
	"github.com/go-rel/rel"
)

const version = "0.0.1"

type Mux = chi.Mux

type Repository = rel.Repository
type Adapter = rel.Adapter

var ContextPath = "/"

var Resources embed.FS

type Enki struct {
	AppName  string
	Env      string
	trunk    *Mux
	Routes   *Mux
	DBConfig database.Config
	DB       Repository
	Shutdown []func() error
}

func New(name string) Enki {
	app := Enki{}

	app.AppName = name

	// Initialize environment
	app.Env = app.fetchEnvironment()

	// Initialize routes
	app.trunk, app.Routes = app.appRoutes()

	return app
}

func (enki *Enki) fetchEnvironment() string {
	env := spt.FetchEnv("ENKI_ENV", "development")

	switch env {
	case "development":
		// do something
	case "production":
		// do something
	case "test":
		// do something
	default:
		log.Fatalf("Invalid environment '%v'", env)
	}

	return env
}

func (enki *Enki) NewDBConfig() database.EnvConfig {
	blob := dbConfigFile()

	config := database.NewConfig(blob, enki.Env)

	return config
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

	halt := make(chan struct{})

	ctx := context.Background()

	go enki.gracefulShutdown(ctx, server, halt)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-halt
}

func (enki *Enki) gracefulShutdown(ctx context.Context, server *http.Server, halt chan struct{}) {
	sigint := make(chan os.Signal, 1)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	log.Println("shutting down server gracefully")

	// stop receiving any request.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("shutdown error", err)
	}

	// close any other things db, redis, etc.
	for i := range enki.Shutdown {
		enki.Shutdown[i]()
	}

	close(halt)
}

func (enki *Enki) Version() string {
	return version
}

func dbConfigFile() []byte {
	file, err := Resources.ReadFile("config/database.yml")

	if err != nil {
		log.Fatal(err.Error())
	}

	return file
}
