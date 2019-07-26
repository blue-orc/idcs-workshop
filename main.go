package main

import (
	"github.com/joho/godotenv"
	"idcs-workshop/api"
	"idcs-workshop/config"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	config.CheckConfig()
	api.Start()
}
