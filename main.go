package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tempo-nksn/Tempo-Backend/db"
	"github.com/tempo-nksn/Tempo-Backend/server"
)

func main() {
	// Opening the database connection
	database, err := gorm.Open("sqlite3", "db/seed/database.db")
	if err != nil {
		panic("failed to establish database connection")
	}
	defer database.Close()
	db.Init(database)

	// Router creation
	router := server.CreateRouter()
	server.StartServer(router)
}
