package database

import "Parser/config"

type Goods struct {
	ID     int    `gorm:"id"`
	Name   string `gorm:"name"`
	URL    string `gorm:"url"`
	URLImg string `gorm:"url_img"`
	Price  string `gorm:"price"`
}

type DBIFace interface {
	Init(*config.DbSettings) error
	AddNewProduct(*Goods) error
}
