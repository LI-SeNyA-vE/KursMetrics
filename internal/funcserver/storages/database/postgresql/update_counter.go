package postgresql

func (d *DataBase) UpdateCounter(name string, value int64) int64 {
	_, err := d.db.Exec(queryUpdateCounter, name, value)
	if err != nil {
		d.log.Printf("ошибка обновления counter: %v", err)
	}

	result, err := d.GetCounter(name)
	if err != nil {
		return 0
	}

	return *result
}
