package server

import (
	"net/http"
)

func (c *Config) handlerHealthz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := c.RawDB.PingContext(ctx); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
