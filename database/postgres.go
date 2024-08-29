package database

import (
	"fmt"
)

func UrlForPostgres(conf Config) string {
	// postgres://username:password@localhost:5432/my_db?sslmode=disable
	format := "%s://%s:%s@%s:%s/%s?sslmode=disable"

	url := fmt.Sprintf(format, conf.Adapter, conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)

	return url
}
