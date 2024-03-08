package main

import (
	"log"

	"portifolio/pkg/handler"
)

func main() {
	log.Println("Starting the application...")
	handler.InitHandlers()
}
