package models

import "gorm.io/gorm"

func AutoMigration(db *gorm.DB) error {
	err := db.AutoMigrate(&View{}, &ViewComponent{}, &DataSeries{}, &DataPoint{}, &DataSource{}, &User{})
	if err != nil {
		return err
	}

	db.Exec(`SELECT create_hypertable('data_series', 'timestamp');`)
	db.Exec(`CREATE INDEX IF NOT EXISTS ix_data_series_datapointid_timestamp ON data_series (data_point_id, timestamp DESC);

`)
	return nil
}
