package main

import (
	"Parser/config"
	"Parser/database"
	"Parser/ozonparser"
	"log"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	settings, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := database.NewDBService()
	err = db.Init(settings.DB)
	if err != nil {
		log.Fatal(err)
	}

	parser := ozonparser.NewParser()
	err = parser.Init(settings.Parser, db)
}
