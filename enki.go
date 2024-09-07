package enki

import (
	"embed"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/JesseChavez/enki/database"
	"github.com/JesseChavez/enki/helper"
	"github.com/JesseChavez/enki/renderer"
	"github.com/JesseChavez/enki/session"
	"github.com/JesseChavez/spt"
	"github.com/go-chi/chi/v5"
	"github.com/go-rel/mssql"
	"github.com/go-rel/postgres"
	"github.com/go-rel/rel"
	"github.com/gorilla/securecookie"

	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
)

const version = "0.0.1"

type Mux = chi.Mux

type Router = chi.Router

type Repository = rel.Repository

type IRenderer interface {
	Render(w http.ResponseWriter, r *http.Request, view string, data any)
}

type IHelper interface {
	RoutePath(string) string
	AssetPath(string) string
	URLParam(*http.Request, string) string
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
	Renderer     IRenderer
	Helper       IHelper
	SessionStore ISessionStore
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

// private variables
var webPort string
var timeZone string
var contextPath string
var sessionKey string
var sessionMaxAge int
var rootPath string

var secretAuthKey string
var secretEncrKey string

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
	rootPath = BaseDir

	secretAuthKey = SecretAuthKey
	secretEncrKey = SecretEncrKey

	return app
}

func (ek *Enki) Version() string {
	return version
}

func (ek *Enki) InitWebApplication(contextMux *Mux) {
	// initialize session store
	ek.SessionStore = session.New(sessionKey, sessionMaxAge, secretAuthKey, secretEncrKey)

	// init db
	intializeDatabase(ek)

	// init renderers
	ek.Renderer = renderer.New(ek.Env, contextPath, rootPath, Resources)

	// init renderers
	ek.Helper = helper.New(ek.Env, contextPath)

	// add shutdown server endpoint
	ek.Routes.Get("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Sutdown request")
		w.Write([]byte("OK"))
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
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
