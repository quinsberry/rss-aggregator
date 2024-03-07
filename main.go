package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Init")
	godotenv.Load()

	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT is not found in the environment")
	}

	fmt.Printf("Port: %s", portStr)
}
