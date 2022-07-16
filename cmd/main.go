package main

import (
	"Parser/config"
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

}
