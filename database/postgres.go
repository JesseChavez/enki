package database

import (
	"fmt"
	"net/url"
)

func UrlForPostgres(conf Config) (string, string) {
	// postgres://username:password@localhost:5432/my_db?sslmode=disable
	format := "%s://%s:%s@%s:%s/%s?sslmode=%s"

	password := url.QueryEscape(conf.Password)

	realUrl := fmt.Sprintf(format, conf.Adapter, conf.Username, password, conf.Host, conf.Port, conf.Database, conf.Sslmode)

	fakeUrl := fmt.Sprintf(format, conf.Adapter, "[∗∗∗∗]", "[∗∗∗∗]", conf.Host, conf.Port, conf.Database, conf.Sslmode)

	return realUrl, fakeUrl
}
