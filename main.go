package main

import "github.com/Tempo-Backend/server"

func main() {
	router := server.CreateRouter()
	server.StartServer(router)
}
