package main

import "Tempo-Backend/server"

func main() {
	router := server.CreateRouter()
	server.StartServer(router)
}
