package server

import "net/http"

func (c *Config) handlerFetchAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assetName := r.URL.Query().Get("asset_name")
	projectName := r.URL.Query().Get("project_name")

	if assetName == "" || projectName == "" {
		respondError(w, http.StatusBadRequest, "missing asset_name or project_name")
		return
	}
}
