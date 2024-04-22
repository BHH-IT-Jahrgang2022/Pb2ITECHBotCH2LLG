package main

import (
	client "client/client"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".ENV problem: ", err)
	}
	client.Initializer()

}
