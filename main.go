package main

import (
	"file-sync/server"
	"log"
)

func main() {
	log.Println("file sync service launching...")
	server.ServeGin()
}
