package handlers

import (
	"database/sql"
	"net/http"
)

// Ping Кидает запрос в базу, для прорки её наличия
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", h.cfg.FlagDatabaseDsn)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
