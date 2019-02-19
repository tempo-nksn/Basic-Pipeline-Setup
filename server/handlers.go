package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"github.com/tempo-nksn/Tempo-Backend/models"
	"googlemaps.github.io/maps"
)

type places struct {
	origin      string
	destination string
}

// RouteInfoStrcut is for getting all information between source and destination
type RouteInfoStrcut struct {
	Allpolyline []string
	Duration    float64
	Distance    int
	Fare        int
	ETA         int
}

func hello(c *gin.Context) {
	c.String(200, "Hello User, your taxi is booked")
}

func dummyPost(c *gin.Context) {
	var str struct {
		Name string `json:"name"`
	}
	c.BindJSON(&str)
	fmt.Printf("I am posting to ios")
	c.JSON(200, str)
}

// getNearByTaxis calculates 3 random taxi location neat the user
// we are using rand func
// The randomness is [-0.002500, 0.002500] degrees for both latitude and longitude
// In future it will respond with the location of the nearby taxis
func getNearByTaxis(c *gin.Context) {
	var userLocation models.Location

	r := c.Request
	m, _ := url.ParseQuery(r.URL.RawQuery)

	if _, ok := m["latitude"]; !ok {
		c.JSON(400, "User Location Missing!!!!!")
		return
	}
	if _, ok := m["longitude"]; !ok {
		c.JSON(400, "User Location Missing!!!!!")
		return
	}
	userLocation.Latitude = m["latitude"][0]
	userLocation.Longitude = m["longitude"][0]
	latitude, _ := strconv.ParseFloat(userLocation.Latitude, 64)
	longitude, _ := strconv.ParseFloat(userLocation.Longitude, 64)

	fmt.Printf("%f", latitude)
	fmt.Printf("%f", longitude)

	//c.JSON(200, userLocation)
	var numberOfTaxis = 3
	var taxiLocation = make([]models.Location, numberOfTaxis)
	for count := 0; count < numberOfTaxis; count++ {
		taxiLocation[count].Latitude = fmt.Sprintf("%0.6f", (latitude + float64(rand.Intn(5000)-2500)/100000))
		taxiLocation[count].Longitude = fmt.Sprintf("%0.6f", (longitude + float64(rand.Intn(5000)-2500)/100000))
	}

	c.JSON(200, taxiLocation)

}

func getPolyLine(c *gin.Context) {

	r := c.Request
	m, _ := url.ParseQuery(r.URL.RawQuery)

	if _, ok := m["src"]; !ok {
		c.JSON(400, "user source Missing!!!!!")
		return
	}
	if _, ok := m["dest"]; !ok {
		c.JSON(400, "User destination Missing!!!!!")
		return
	}

	srcAndDest := places{origin: m["src"][0], destination: m["dest"][0]}
	allPolyline := getRoute(srcAndDest)
	c.JSON(200, allPolyline)
}

// func confirmBooking(c *gin.Context) {
// 	// 1) Get the user id, taxi id, route id from context
// 	q := c.Request.URL.Query()
// 	userId := q["userid"][0]
// 	taxiId := q["taxiid"][0]
// 	routeId := q["routeid"][0]
// 	fmt.Println(userId, taxiId, routeId)
// 	// 2) In the database, in Bookings table, store all the values

// 	// 3) Respond to the server saying booking is done, return a string saying "booking done"
// 	c.String(200, "Booking done")
// }

// all fuction supporting the API fucntions

func getRoute(obj places) RouteInfoStrcut {

	// googleKey := os.Getenv("MAPS_KEY")
	// c, err := maps.NewClient(maps.WithAPIKey(googleKey))

	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBZgBxptkNs1pSNUsy-QOtCQRimm4JCNDE"))

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DirectionsRequest{
		Origin:      obj.origin,
		Destination: obj.destination,
	}
	r.Mode = maps.TravelModeDriving
	r.Units = maps.UnitsMetric
	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	var maplegs *maps.Leg
	maplegs = route[0].Legs[0]
	pretty.Println(maplegs)
	var allPolyline []string
	var polyline string

	for i := 0; i < len(maplegs.Steps); i++ {

		polyline = maplegs.Steps[i].Polyline.Points
		allPolyline = append(allPolyline, polyline)
	}
	routeInfo := RouteInfoStrcut{Allpolyline: allPolyline, Distance: route[0].Legs[0].Distance.Meters, Duration: route[0].Legs[0].Duration.Minutes(), Fare: getfare(route[0].Legs[0].Distance.Meters), ETA: getETA(route[0].Legs[0].Duration.Minutes())}

	return routeInfo

}

func getfare(distance int) int {
	basePricePerMeter := 3
	return int(distance/1000) * basePricePerMeter //devided by thousand because distance is in meters

}

func getETA(duration float64) int {
	minWaitingTime, maxWaitingTime := 5, 15

	waitingTime := int(0.05 * duration)
	if (waitingTime < maxWaitingTime) && (waitingTime > minWaitingTime) {
		return waitingTime
	} else if waitingTime > maxWaitingTime {
		return maxWaitingTime

	} else {
		return minWaitingTime
	}
}
