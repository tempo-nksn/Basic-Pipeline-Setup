package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tempo-nksn/Tempo-Backend/db"
	"github.com/tempo-nksn/Tempo-Backend/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	var taxis []models.Taxi
	var riders []models.Rider
	var drivers []models.Driver
  var routes []models.Route
	getData("data/taxis.json", &taxis)
	getData("data/riders.json", &riders)
	getData("data/drivers.json", &drivers)
	getData("data/routes.json", &routes)

	DATABASE := os.Getenv("DB_DRIVER")
	databaseURL := os.Getenv("DATABASE_URL")
	if DATABASE == "" && databaseURL == "" {
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		DATABASE = os.Getenv("DB_DRIVER")
		databaseURL = os.Getenv("DATABASE_URL")
	}
	// team Roses postgres URL: postgres://fpfujvlpxoelcm:93ec17a8d9323c29e05b569ea2bd77fa2d7dc96564d2f5b6eaa521807bf8b787@ec2-54-243-128-95.compute-1.amazonaws.com:5432/d34p830c249n99
	database, err := gorm.Open(DATABASE, databaseURL)
	if err != nil {
		panic("failed to establish database connection")
	}
	defer database.Close()
	db.Init(database)
	db.DB.Unscoped().Delete(&taxis)
	// db.DB.Delete(&taxis)
	for _, taxi := range taxis {
		db.DB.Create(&taxi)
	}
	for _, rider := range riders {
		db.DB.Create(&rider)
	}
	for _, driver := range drivers {
		db.DB.Create(&driver)
  }
	for _, route := range routes {
		db.DB.Create(&route)
	}
	log.Println("Seed data created!! Now you move it!")
}

func getData(fileName string, v interface{}) {
	file, _ := os.Open(fileName)
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, v)
}
