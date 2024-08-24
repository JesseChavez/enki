package enki

import (
	"embed"
	"log"
	"net/http"

	"github.com/JesseChavez/enki/database"
	"github.com/JesseChavez/enki/templating"
	"github.com/JesseChavez/spt"
	"github.com/go-chi/chi/v5"
	"github.com/go-rel/mssql"
	"github.com/go-rel/rel"

	_ "github.com/microsoft/go-mssqldb"
)

const version = "0.0.1"

type Mux = chi.Mux

type Repository = rel.Repository

type Renderer interface {
	Render (w http.ResponseWriter, r *http.Request, view string, data any)
}

type Enki struct {
	AppName  string
	Env      string
	Routes   *Mux
	DBConfig database.Config
	DB       Repository
	Render   Renderer
}

var BaseDir string

var Resources embed.FS

var ContextPath = "/"

var WebPort  = "3000"

var TimeZone = "UTC"


// private variables
var webPort     string
var timeZone    string
var contextPath string
var rootPath    string

func New(name string) Enki {
	app := Enki{}

	app.AppName = name

	webPort     = WebPort
	timeZone    = TimeZone
	contextPath = ContextPath
	rootPath    = BaseDir

	return app
}

func (ek *Enki) Version() string {
	return version
}

func (ek *Enki) InitWebApplication (contextMux *Mux) {
	// Initialize environment
	ek.Env = ek.fetchEnvironment()

	// init db
	intializeDatabase(ek)

	// init renderers
	ek.Render = templating.New( ek.Env, rootPath, Resources)
}

func (ek *Enki) InitJobApplication () {
	// Initialize environment
	ek.Env = ek.fetchEnvironment()
}

func (ek *Enki) InitDbMigration () {
	// Initialize environment
	ek.Env = ek.fetchEnvironment()
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

func intializeDatabase(ek *Enki) {
	config := ek.NewDBConfig()

	adapterName := config.Current.Adapter

	url := config.Current.GetUrl()

	var adapter rel.Adapter
	var err error

	switch adapterName {
	case "sqlserver":
		adapter, err = mssql.Open(url)
	case "postgres":
		// adapter, err = postgres.Open(url)
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
	blob := dbConfigFile()

	config := database.NewConfig(blob, ek.Env)

	return config
}

func dbConfigFile() []byte {
	file, err := Resources.ReadFile("config/database.yml")

	if err != nil {
		log.Fatal(err.Error())
	}

	return file
}
