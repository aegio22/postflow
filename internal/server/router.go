package server

import (
	"net/http"

	"github.com/aegio22/postflow/internal/routes"
)

func (c *Config) NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("POST "+routes.SignUp, c.handlerSignUp)
	r.HandleFunc("POST "+routes.Login, c.handlerLogin)

	return r
}
