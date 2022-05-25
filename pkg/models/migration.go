package models

import "gorm.io/gorm"

func AutoMigration(db *gorm.DB) error {
	return db.AutoMigrate(&DataSource{}, &User{})
}
