package main

import (
	"dexbg/api"
	"dexbg/storage"
	"fmt"
	"log"
)

func main() {
	var listenAddr string = "3000"

	store := storage.NewMongoStorage()
	store.ConnectMongoDB()
	defer store.DisconnectClient()

	server := api.NewServer(listenAddr, store)
	fmt.Println("Server running on port: " + listenAddr)
	log.Fatal(server.Start())
}
