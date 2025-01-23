package postgresql

func (d *DataBase) UpdateGauge(name string, value float64) float64 {
	_, err := d.db.Exec(queryUpdateGauge, name, value)
	if err != nil {
		d.log.Printf("ошибка обновления gauge: %v", err)
	}
	return value
}
