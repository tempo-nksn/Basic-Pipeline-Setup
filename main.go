package main

import "github.com/tempo-nksn/Tempo-Backend/server"

func main() {
	router := server.CreateRouter()
	server.StartServer(router)
}
