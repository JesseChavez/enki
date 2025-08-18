package bouncer

import (
	"encoding/json"
	"net/http"
	"net/url"
)


type CookieStore struct {
	BaseOptions *Options
}


func NewCookieStore() *CookieStore {
	store := CookieStore{}

	return &store
}


func (cs *CookieStore) Get(r *http.Request, name string) (*Session, error) {
	session := Session{
		name: name,
		store: cs,
	}

	opts := *cs.BaseOptions
	session.Options = &opts
	session.Values = make(map[string]any)

	cookie, err := r.Cookie(name)

	if err != nil {
		return  &session, err
	}

	err = cs.Unpack(cookie.Value, &session.Values)

	return  &session, err
}

func (cs *CookieStore) Save(r *http.Request, w http.ResponseWriter, s *Session) error {
	// err will be nil for cookie store but other stores may reture error
	// e.g  redis or database connection error
	var err error


	cookie := cs.NewCookie(s.name, s.Options)

	cookie.Value, err = cs.Pack(s.Values)

	http.SetCookie(w, cookie)
	return err
}

func (cs *CookieStore) Delete(w http.ResponseWriter, s *Session) error {
	// err will be nil for cookie store but other stores may reture error
	// e.g  redis or database connection error
	var err error

	value := "{}"
	
	cookie := cs.NewCookie(s.name, s.Options)

	cookie.Value  = value
	cookie.MaxAge = -1

	http.SetCookie(w, cookie)
	return err
}

func (cs *CookieStore) NewCookie(name string, options *Options) *http.Cookie {
	cookie := http.Cookie{}

	cookie.Name        = name

	// cookie options
	cookie.Path        = options.Path
	cookie.Domain      = options.Domain
	cookie.MaxAge      = options.MaxAge
	cookie.Secure      = options.Secure
	cookie.HttpOnly    = options.HttpOnly
	cookie.Partitioned = options.Partitioned
	cookie.SameSite    = options.SameSite

	return &cookie
}

func (cs *CookieStore) Unpack(escapedData string, decodedData any) error {
	var err error

	encodedData, err := url.QueryUnescape(escapedData)

	err = json.Unmarshal([]byte(encodedData), decodedData)

	return err
}

func (cs *CookieStore) Pack(decodedData any) (string, error) {
	var err error

	encodedData, err := json.Marshal(decodedData)

	if err != nil {
		return "", err
	}

	escapedData := url.QueryEscape(string(encodedData))

	return escapedData, err
}
