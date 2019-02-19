package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/tempo-nksn/Tempo-Backend/db"
	"github.com/tempo-nksn/Tempo-Backend/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var taxis []models.Taxi

	getData("data/taxis.json", &taxis)
	database, err := gorm.Open("postgres", "postgres://fpfujvlpxoelcm:93ec17a8d9323c29e05b569ea2bd77fa2d7dc96564d2f5b6eaa521807bf8b787@ec2-54-243-128-95.compute-1.amazonaws.com:5432/d34p830c249n99")
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
