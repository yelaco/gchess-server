package main

import (
	"github.com/joho/godotenv"
)

func main() {
	// Look for a file name .env in current directory to load environment variables
	godotenv.Load()

}
