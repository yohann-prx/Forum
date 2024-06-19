package main

import (
	"SPORTALK/internal/app/server"
	"log"
)

func main() {
	// Create a new configuration instance
	// Read configuration from file
	config := server.NewConfig()
	if err := config.ReadConfig(); err != nil {
		log.Fatal(err)
	}

	// Start the server with the obtained configuration
	log.Fatal(server.Start(config))
}
