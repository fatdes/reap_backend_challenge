package restutil

import (
	"net/http"
)

type Error struct {
	StatusCode   int `json:"-"`
	ErrorMessage string
}

func (error *Error) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(error.StatusCode)
	return nil
}
