package models

import "gorm.io/gorm"

func AutoMigration(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &View{}, &ViewComponent{}, &DataSource{}, &DataPoint{}, &DataSeries{})
	if err != nil {
		return err
	}

	db.Exec(`CREATE INDEX IF NOT EXISTS ix_data_series_datapointid_timestamp ON data_series (data_point_id, timestamp DESC);`)

	err = CreateHyperTable(db, "data_series", "timestamp")

	if err != nil {
		return err
	}

	return nil
}

func CreateHyperTable(db *gorm.DB, table, column string) error {

	var exists bool

	err := db.Raw(`SELECT
    CASE WHEN EXISTS
    (
        select * from _timescaledb_catalog.hypertable where table_name = ?
    )
    THEN true
    ELSE false
END`, table).Scan(&exists).Error

	if err != nil {
		return err
	}

	if !exists {
		return db.Exec(`SELECT create_hypertable(?, ?);`, table, column).Error
	}

	return nil
}
