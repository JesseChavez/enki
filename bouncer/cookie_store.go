package bouncer

import (
	"net/http"
	"net/url"
)


type CookieStore struct {
	BaseOptions *Options
	Transcoder  *Transcoder
}


func NewCookieStore(secret string, salt string, maxAge int) *CookieStore {
	store := CookieStore{}

	transcoder := new(Transcoder).Init(secret, salt, maxAge)

	store.Transcoder = transcoder

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

	err = cs.Decode(cookie.Value, &session.Values)

	return  &session, err
}

func (cs *CookieStore) Save(r *http.Request, w http.ResponseWriter, s *Session) error {
	// err will be nil for cookie store but other stores may reture error
	// e.g  redis or database connection error
	var err error


	cookie := cs.NewCookie(s.name, s.Options)

	cookie.Value, err = cs.Encode(s.Values)

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

func (cs *CookieStore) Decode(escapedValue string, decodedValue any) error {
	var err error

	encodedValue, err := url.QueryUnescape(escapedValue)

	err = cs.Transcoder.Decode(encodedValue, decodedValue)

	return err
}

func (cs *CookieStore) Encode(decodedValue any) (string, error) {
	var err error

	encodedValue, err := cs.Transcoder.Encode(decodedValue)

	if err != nil {
		return "", err
	}

	escapedValue := url.QueryEscape(encodedValue)

	return escapedValue, err
}
