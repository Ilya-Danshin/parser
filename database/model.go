package database

import "Parser/config"

type DBIFace interface {
	Init(config.DbSettings) error
}
