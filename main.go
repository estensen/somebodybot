package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	token := os.Getenv("SLACKTOKEN")
	if token == "" {
		log.Fatal("missing Slack token env var")
	}	
	fmt.Println(token)
}

