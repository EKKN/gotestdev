package main

import (
	"log"
	"os"

	"github.com/EKKN/gotestdev/packages/api"
	"github.com/EKKN/gotestdev/packages/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config.ConnectDatabase()

	APP_PORT := os.Getenv("APP_PORT")
	newApi := api.NewAPI(APP_PORT)
	newApi.Run()
}
