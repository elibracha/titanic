package main

import (
	"log"
	"titanic-api/internal"
)

// @title           Titanic API
// @version         1.0
// @description     This is API provide multiple functionality endpoints over titanic dataset
// @contact.name    Eli Bracha
// @BasePath       /api/v1
func main() {
	server, err := internal.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	server.Start()
}
