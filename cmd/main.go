package main

import (
	"log"

	"Parser/config"
	"Parser/database"
	"Parser/ozonparser"
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

	err = parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
}
