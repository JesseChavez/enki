package database

import (
	"embed"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type EnvConfig struct {
	Development Config
	Test        Config
	Production  Config
	Current     Config
}

type Config struct {
	Adapter  string `yaml:"adapter,omitempty"`
	Encoding string `yaml:"encoding,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Port     string `yaml:"port,omitempty"`
	Database string `yaml:"database,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Sslmode  string `yaml:"sslmode,omitempty"`
}

func NewConfig(file []byte, env string) EnvConfig {
	// var config map[string]map[string]string
	config := EnvConfig{}
	var err error

	err = yaml.Unmarshal(file, &config)

	if err != nil {
		log.Fatal(err.Error())
	}

	// fmt.Printf("db config: %+v\n", config)

	curr := config.GetEnv(env)

	config.Current = curr

	return config
}

func (conf *Config) GetUrl() string {
	var url string

	adapter := conf.Adapter

	switch adapter {
	case "sqlserver":
		url = UrlForMssql(*conf)
	case "postgres":
		url = UrlForPostgres(*conf)
	default:
		log.Fatalf("Invalid database adapter '%v'", adapter)
	}

	return url
}

func (conf *EnvConfig) GetEnv(env string) Config {
	var params Config

	switch env {
	case "development":
		params = conf.Development
	case "production":
		params = conf.Production
	case "test":
		params = conf.Test
	default:
		log.Fatalf("Invalid application environment '%v'", env)
	}

	return params
}

func ConfigFile(appName string, env string, rootPath string, resources embed.FS) []byte {
	if env == "production" {
		fileName := "database_" + appName + ".yml"

		systemFile := "/var/local/config/" + fileName

		file, err := os.ReadFile(systemFile)

		if err == nil {
			log.Println("loading file:", systemFile)
			expandedFile := os.ExpandEnv(string(file))
			return []byte(expandedFile)
		}
		
		log.Println("DB system file not found:", systemFile)

		systemWorkingDirFile := rootPath + "/" + fileName

		file, err = os.ReadFile(systemWorkingDirFile)

		if err == nil {
			log.Println("loading file:", systemFile)
			expandedFile := os.ExpandEnv(string(file))
			return []byte(expandedFile)
		}

		log.Println("DB system file not found:", systemWorkingDirFile)

		globalFile := rootPath + "/database.yml"

		file, err = os.ReadFile(globalFile)

		if err == nil {
			log.Println("loading file:", globalFile)
			expandedFile := os.ExpandEnv(string(file))
			return []byte(expandedFile)
		}

		log.Println("DB global file not found:", globalFile)
		log.Println("Trying embeded file")
	}

	file, err := resources.ReadFile("config/database.yml")

	log.Println("Project root:", rootPath)

	if err != nil {
		log.Fatal(err.Error())
	}

	expandedFile := os.ExpandEnv(string(file))
	return []byte(expandedFile)
}
