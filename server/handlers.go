package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tempo-nksn/Tempo-Backend/models"
	"math/rand"
	"net/url"
	"strconv"
)

//func hello(c *gin.Context) {
//	c.String(200, "Hello User, your taxi is booked")
//}

func dummyPost(c *gin.Context) {
	var str struct {
		Name string `json:"name"`
	}
	c.BindJSON(&str)
	fmt.Printf("I am posting to ios")
	c.JSON(200, str)
}

// getTaxi calculates 3 random taxi location neat the user
// In future it will respond with the location of the nearby taxis
func getNearByTaxi(c *gin.Context) {
	var userLocation models.Location

	r := c.Request
	m,_ := url.ParseQuery(r.URL.RawQuery)

	userLocation.Latitude = m["latitude"][0]
	userLocation.Longitude = m["longitude"][0]
	latitude,_ := strconv.ParseFloat(userLocation.Latitude, 64)
	longitude,_ := strconv.ParseFloat(userLocation.Longitude, 64)

	fmt.Printf("%f",latitude)
	fmt.Printf("%f",longitude)

	//c.JSON(200, userLocation)

	var taxiLocation [3]models.Location
	for count := 0; count<3; count++ {
		taxiLocation[count].Latitude = fmt.Sprintf("%0.6f",(latitude + float64(rand.Intn(5000) - 2500)/100000))
		taxiLocation[count].Longitude = fmt.Sprintf("%0.6f",(longitude + float64(rand.Intn(5000) - 2500)/100000))
	}

	c.JSON(200, taxiLocation)

}

