package database

import (
	"fmt"
)

func UrlForMssql(conf Config) string {
	// "sqlserver://username:password@localhost:1433?database=my_db"
	format := "%s://%s:%s@%s:1433?database=%s"

	url := fmt.Sprintf(format, conf.Adapter, conf.Username, conf.Password, conf.Host, conf.Database)

	return url
}
