package server

import "net/http"

func (c *Config) NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	// example: r.HandleFunc("POST /users", c.CreateUser)

	return r
}
