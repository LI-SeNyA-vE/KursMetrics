/*
Package handlers содержит функции-обработчики (Handler),
управляющие различными HTTP-запросами на сервер.
Функция Ping проверяет доступность базы данных посредством
отправки запроса db.Ping() к PostgreSQL.
*/
package handlers

import (
	"database/sql"
	"net/http"
)

// Ping осуществляет проверку соединения с базой данных по строке подключения,
// указанной в h.cfg.FlagDatabaseDsn. Если соединение успешно устанавливается
// (db.Ping() без ошибок), возвращается статус 200 (OK).
// В противном случае возвращаются коды 404 или 500, указывающие на проблему
// при открытии или пинге базы данных.
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
