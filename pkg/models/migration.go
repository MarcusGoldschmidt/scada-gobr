package models

import "gorm.io/gorm"

func AutoMigration(db *gorm.DB) error {
	err := db.AutoMigrate(&DataSeries{}, &DataPoint{}, &DataSource{}, &User{})
	if err != nil {
		return err
	}

	db.Exec(`SELECT create_hypertable('data_series', 'time_stamp');`)
	db.Exec(`CREATE INDEX IF NOT EXISTS ix_data_series_datapointid_time_stamp ON data_series (data_point_id, time_stamp DESC);

`)
	return nil
}
