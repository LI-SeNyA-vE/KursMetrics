package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func CrereateDB(db *sql.DB, createTableSQL string) {
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Println("Ошибка при создании таблицы: %v", err)
	}
	fmt.Println("Таблица 'Metric' успешно создана или уже была.")
}
