package main

import (
	"log"

	"github.com/e0m-ru/fileserver/internal/server"
)

func main() {
	err := server.StartServer()
	if err != nil {
		log.Print(err)
	}
}
