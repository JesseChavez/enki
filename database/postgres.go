package database

import (
	"fmt"
)

func UrlForPostgres(conf Config) string {
	// postgres://username:password@localhost:5432/my_db
	format := "%s://%s:%s@%s:5432/%s"

	url := fmt.Sprintf(format, conf.Adapter, conf.Username, conf.Password, conf.Host, conf.Database)

	return url
}
