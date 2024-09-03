package session

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)


type SessionStore struct {
	Name string
	Store *sessions.CookieStore
}


type Session struct {
	req *http.Request
	val *sessions.Session
}

func New(name string, authKey string, encrKey string) *SessionStore {
	store := SessionStore{Name: name}

	// authKeyOne := securecookie.GenerateRandomKey(64)
	// encrKeyOne := securecookie.GenerateRandomKey(32)
	authKeyOne := []byte(authKey)
	encrKeyOne := []byte(encrKey)

	cstore := sessions.NewCookieStore(authKeyOne, encrKeyOne)

	cstore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	store.Store = cstore

	return &store
}

func (ss *SessionStore) GetSession(r *http.Request) *Session {
	session, err := ss.Store.Get(r, ss.Name)
	
	if err != nil {
		log.Println("Failure on retrieving session", err)
	}
	
	return &Session{val: session, req: r}
}

func (s *Session) Get(field string) any {
	value := s.val.Values[field]

	return value
}

func (s *Session) Set(field string, value any) {	
	s.val.Values[field] = value
}

// Save the current session.
func (s *Session) Save(w http.ResponseWriter) error {
	return s.val.Save(s.req, w)
}

func (s *Session) Delete(w http.ResponseWriter) error {
	s.val.Options.MaxAge = -1

	return s.val.Save(s.req, w)
}
