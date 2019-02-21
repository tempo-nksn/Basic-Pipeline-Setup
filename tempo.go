package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tempo-nksn/Tempo-Backend/db"
	"github.com/tempo-nksn/Tempo-Backend/server"
)

func main() {
	// Connecting to the database
	database, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		panic("failed to establish database connection")
	}
	defer database.Close()
	db.Init(database)

	// Creating a router
	router := server.CreateRouter(database)
	server.StartServer(router)
}
