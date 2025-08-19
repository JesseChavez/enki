package enki

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JesseChavez/enki/commands"
	"github.com/JesseChavez/enki/database"
	"github.com/JesseChavez/enki/logger"
	"github.com/JesseChavez/enki/bouncer"
	"github.com/JesseChavez/enki/view"
	"github.com/JesseChavez/spt"
	"github.com/go-chi/chi/v5"
	"github.com/go-rel/mssql"
	"github.com/go-rel/postgres"
	"github.com/go-rel/rel"

	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
)

const version = "0.4.4"

type Mux = chi.Mux

type Router = chi.Router

type Repository = rel.Repository

type ActionView = view.ActionView

type ILogger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

type IViewSupport interface {
	RoutePath(string) string
	AssetPath(string) string
	URLParam(*http.Request, string) string
	Render(w http.ResponseWriter, status int, view *ActionView)
	RenderHTML(w http.ResponseWriter, status int, view *ActionView)
	RenderXML(w http.ResponseWriter, status int, view *ActionView)
}

type ISessionManager interface {
	GetSession(*http.Request) *bouncer.Session
}

type Enki struct {
	AppName      string
	Env          string
	Routes       *Mux
	Queues       *Queues
	DBConfig     database.Config
	DB           Repository
	Logger       ILogger
	SessionManager ISessionManager
	ViewSupport    IViewSupport
}


var BaseDir string

var Resources embed.FS

var ContextPath = "/"

var SessionKey = "_enki_session"

var SessionMaxAge = 30

var WebPort = "3000"

var TimeZone = "UTC"

var SecretKeyBase = "secret-key-base"


var AuthenticatedEncryptedCookieSalt = "authenticated encrypted cookie"

var API = false

var CSR = false

// private variables
var webPort string
var timeZone string
var contextPath string
var sessionKey string
var sessionMaxAge int
var rootPath string

var secretKeyBase string
var authenticatedEncryptedCookieSalt string
var api bool
var csr bool

func New(name string) Enki {
	app := Enki{}

	// Initialize environment before anything else
	app.Env = app.fetchEnvironment()

	app.AppName = name

	webPort = WebPort
	timeZone = TimeZone
	contextPath = ContextPath
	sessionKey = SessionKey
	sessionMaxAge = SessionMaxAge

	rootPath = workingDir()

	secretKeyBase = SecretKeyBase
	authenticatedEncryptedCookieSalt = AuthenticatedEncryptedCookieSalt

	api = API
	csr = CSR

	return app
}

func (ek *Enki) Version() string {
	return version
}

func (ek *Enki) InitWebApplication(contextMux *Mux) {
	// init logger
	ek.Logger = logger.New()

	// initialize session manager
	ek.SessionManager = bouncer.New(
		sessionKey, secretKeyBase, authenticatedEncryptedCookieSalt, sessionMaxAge, false,
	)

	// init db
	intializeDatabase(ek)

	// init support (renderer and helpers)
	ek.ViewSupport = view.New(ek.Env, api, csr, contextPath, rootPath, Resources)

	// Enable static assets server
	ek.staticAssets(contextPath, ek.Routes)

	// add shutdown server endpoint
	ek.Routes.Get("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Sutdown request")
		w.Write([]byte("OK"))


		StatusCode <- "stop"
	})
}

func (ek *Enki) InitJobApplication() {
	// init logger
	ek.Logger = logger.New()

	// init db
	intializeDatabase(ek)
}

func (ek *Enki) ExecuteCommand(command []string) {
	runner := commands.Runner{
		Env: ek.Env,
		Command: command,
	}

	runner.Perform()
}

func (enki *Enki) fetchEnvironment() string {
	env := spt.FetchEnv("APP_ENV", "development")

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

func intializeDatabase(ek *Enki) {
	config := ek.NewDBConfig()

	adapterName := config.Current.Adapter

	url := config.Current.GetUrl()

	log.Println("DB connection url:", url)

	var adapter rel.Adapter
	var err error

	switch adapterName {
	case "sqlserver":
		adapter, err = mssql.Open(url)
	case "postgres":
		adapter, err = postgres.Open(url)
	case "sqlite3":
	// adapter, err = sqlite3.Open(url)
	default:
		log.Fatalf("Invalid adapter '%v'", adapterName)
	}

	if err != nil {
		panic(err)
	}

	// Add to shutdown list.
	Shutdown = append(Shutdown, adapter.Close)

	ek.DB = rel.New(adapter)
	ek.DB.Instrumentation(func(ctx context.Context, op string, message string, args ...any) func(err error) {
		// no op for rel functions.
		if strings.HasPrefix(op, "rel-") {
			return func(error) {}
		}

		t := time.Now()

		return func(err error) {
			duration := time.Since(t)
			stats := "[duration: " + fmt.Sprint(duration) + " op: " + op + "]"

			if err != nil {
				ek.Logger.Error(message, "stat", stats, "err", err)
			} else {
				ek.Logger.Info(message, "stat", stats)
			}
		}
	})

	// repo.Ping(context.TODO())
}

func (ek *Enki) NewDBConfig() database.EnvConfig {
	blob := database.ConfigFile(ek.Env, rootPath, Resources)

	config := database.NewConfig(blob, ek.Env)

	return config
}

func workingDir() string {
	exec, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(exec)

	// detect `go run` otherwise is deployable artifact
	if !strings.Contains(dir, "go-build") {
		return filepath.Dir(exec)
	}

	path, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	return path
}
