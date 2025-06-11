package main

import (
	"fmt"
	"log"

	"github.com/andrii/apz-pzpi-22-2-ieremenko-andrii/Lab2/pzpi-22-2-ieremenko-andrii-lab2/server"
)

func main() {
	// Create and start the server
	s := server.NewServer()

	// Start the server in a goroutine
	go func() {
		fmt.Println("Server is starting on http://localhost:8080")
		fmt.Println("Swagger UI is available at http://localhost:8080/swagger/index.html")
		if err := s.Start(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	select {}
}
