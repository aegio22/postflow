package server

import (
	"net/http"
)

func (c *Config) NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	return r
}
