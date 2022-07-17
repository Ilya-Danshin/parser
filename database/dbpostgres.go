package database

import (
	"fmt"

	"Parser/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBService struct {
	db *gorm.DB
}

func NewDBService() *DBService {
	return &DBService{}
}

func (db *DBService) Init(settings *config.DbSettings) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		settings.Host, settings.User, settings.Password, settings.Name, settings.Port, settings.SSLMode)
	_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.db = _db
	return nil
}

func (db *DBService) AddNewProduct(prod *Goods) error {
	err := db.db.Create(prod).Error
	if err != nil {
		return err
	}

	return nil
}
