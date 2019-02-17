package main

import "github.com/codechaitu/circleci-heroku/server"

func main() {
	router := server.CreateRouter()
	server.StartServer(router)
}