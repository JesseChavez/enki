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

	"github.com/JesseChavez/enki/database"
	"github.com/JesseChavez/enki/logger"
	"github.com/JesseChavez/enki/session"
	"github.com/JesseChavez/enki/view"
	"github.com/JesseChavez/spt"
	"github.com/go-chi/chi/v5"
	"github.com/go-rel/mssql"
	"github.com/go-rel/postgres"
	"github.com/go-rel/rel"
	"github.com/gorilla/securecookie"

	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
)

const version = "0.3.6"

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

type ISessionStore interface {
	GetSession(*http.Request) *session.Session
}

type Enki struct {
	AppName      string
	Env          string
	Routes       *Mux
	DBConfig     database.Config
	DB           Repository
	Logger       ILogger
	SessionStore ISessionStore
	ViewSupport  IViewSupport
}


var BaseDir string

var Resources embed.FS

var ContextPath = "/"

var SessionKey = "_enki_session"

var SessionMaxAge = 30

var WebPort = "3000"

var TimeZone = "UTC"

var SecretAuthKey = string(securecookie.GenerateRandomKey(64))
var SecretEncrKey = string(securecookie.GenerateRandomKey(32))

var API = false

var CSR = false

// private variables
var webPort string
var timeZone string
var contextPath string
var sessionKey string
var sessionMaxAge int
var rootPath string

var secretAuthKey string
var secretEncrKey string
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

	secretAuthKey = SecretAuthKey
	secretEncrKey = SecretEncrKey
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

	// initialize session store
	ek.SessionStore = session.New(sessionKey, sessionMaxAge, secretAuthKey, secretEncrKey)

	// init db
	intializeDatabase(ek)

	// init support (renderer and helpers)
	ek.ViewSupport = view.New(ek.Env, api, csr, contextPath, rootPath, Resources)

	// add shutdown server endpoint
	ek.Routes.Get("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Sutdown request")
		w.Write([]byte("OK"))


		StatusCode <- "stop"
	})
}

func (ek *Enki) InitJobApplication() {
}

func (ek *Enki) InitDbMigration() {
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

func intializeSessionStore(ek *Enki) {
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
	blob := dbConfigFile(ek.Env)

	config := database.NewConfig(blob, ek.Env)

	return config
}

func dbConfigFile(env string) []byte {
	if env == "production" {
		file, err := os.ReadFile(rootPath + "/database.yml")

		if err == nil {
			return file
		}

		log.Println("File not found, trying embeded file")
	}

	file, err := Resources.ReadFile("config/database.yml")

	log.Println("Project root:", rootPath)

	if err != nil {
		log.Fatal(err.Error())
	}

	return file
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
