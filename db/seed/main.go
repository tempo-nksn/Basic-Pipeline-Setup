package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/tempo-nksn/Tempo-Backend/db"
	"github.com/tempo-nksn/Tempo-Backend/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	var taxis []models.Taxi

	getData("data/taxis.json", &taxis)
	database, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic("failed to establish database connection")
	}
	defer database.Close()
	db.Init(database)

	for _, taxi := range taxis {
		db.DB.Save(&taxi)
	}

}

func getData(fileName string, v interface{}) {
	file, _ := os.Open(fileName)
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, v)
}
