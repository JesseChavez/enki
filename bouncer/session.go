package bouncer

import (
	"log"
	"net/http"
)

type Session struct {
	ID string
	// Values contains the user-data for the session.
	Values  map[string]any
	Options *Options
	IsNew   bool
	name    string
	store   Store
}


func (s *Session) Get(field string) any {
	value := s.Values[field]

	return value
}

func (s *Session) Set(field string, value any) {	
	s.Values[field] = value
}

// Save the current session.
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	err := s.store.Save(r, w, s)

	if err != nil {
		log.Println("Sav", err)
	}

	return err
}

func (s *Session) Delete(w http.ResponseWriter) error {
	s.Options.MaxAge = -1

	err := s.store.Delete(w, s)

	if err != nil {
		log.Println("Del", err)
	}

	return err
}
