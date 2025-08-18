package bouncer

import (
	"log"
	"net/http"
)

type Options struct {
       Path        string
       Domain      string
       MaxAge      int
       Secure      bool
       HttpOnly    bool
       Partitioned bool
       SameSite    http.SameSite
}

type Store interface {
	Get(r *http.Request, name string) (*Session, error)
	Save(r *http.Request, w http.ResponseWriter, s *Session) error
	Delete(w http.ResponseWriter, s *Session) error
}

type Manager struct {
	Name string
	Store Store
}

func New(name string, secret string, maxAge int, secure bool) *Manager {
	manager := Manager{Name: name}

	// change to different store here.
	store := NewCookieStore()

	// MaxAge unit is seconds
	store.BaseOptions = &Options{
		Path:     "/",
		MaxAge:   60 * maxAge,
		HttpOnly: true,
		Secure: secure,
	}

	manager.Store = store

	return &manager
}

func (m *Manager) GetSession(r *http.Request) *Session {
	session , err := m.Store.Get(r, m.Name)

	if err != nil {
		log.Println("Error retrivieng session", err)
	}

	return session
}
