package database

import (
	"fmt"
)

func UrlForMssql(conf Config) (string, string) {
	// "sqlserver://username:password@localhost:1433?database=my_db"
	format := "%s://%s:%s@%s:%s?database=%s"

	realUrl := fmt.Sprintf(format, conf.Adapter, conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)

	fakeUrl := fmt.Sprintf(format, conf.Adapter, "[∗∗∗∗]", "[∗∗∗∗]", conf.Host, conf.Port, conf.Database)

	return realUrl, fakeUrl
}
