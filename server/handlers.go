package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kr/pretty"
	"github.com/tempo-nksn/Tempo-Backend/constants"
	"github.com/tempo-nksn/Tempo-Backend/db"
	"github.com/tempo-nksn/Tempo-Backend/models"
	"googlemaps.github.io/maps"
)

// RouteInfoStrcut is for getting all information between source and destination
type RouteInfoStrcut struct {
	Allpolyline []string
	Duration    float64
	Distance    int
	Price       int
	ETA         int
}

func getDB(c *gin.Context) *gorm.DB {
	return c.MustGet(constants.ContextDB).(*gorm.DB)
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

func bookingConfirmation(c *gin.Context) {
	// 1) Get the user id, taxi id, route id from context
	q := c.Request.URL.Query()
	riderId, _ := strconv.Atoi(q["riderid"][0])
	taxiId, _ := strconv.Atoi(q["taxiid"][0])
	routeId, _ := strconv.Atoi(q["routeid"][0])
	fmt.Println(riderId, taxiId, routeId)
	// 2) In the database, in Bookings table, store all the values
	var booking models.Booking
	db.DB.Create(&booking)

	booking.RiderID = uint(riderId)
	booking.TaxiID = uint(taxiId)
	booking.RouteID = uint(routeId)
	// 3) Get the duration, based on Routeid
	var route models.Route
	dbc := getDB(c)
	dbc.Where("id=?", routeId).Find(&route)
	fmt.Println("Duration: ", route.Duration)
	booking.Route.Duration = route.Duration
	db.DB.Save(&booking)
	c.JSON(200, booking)

	// last) Respond to the server saying booking is done, return a string saying "booking done"
	//c.String(200, "Booking done")
}

// all fuction supporting the API fucntions
func getRoute(c *gin.Context) {

	// getting  Source and destination
	request := c.Request
	m, _ := url.ParseQuery(request.URL.RawQuery)
	if _, ok := m["src"]; !ok {
		c.JSON(400, "user source Missing!!!!!")
		return
	}
	if _, ok := m["dest"]; !ok {
		c.JSON(400, "User destination Missing!!!!!")
		return
	}
	origin := m["src"][0]
	destination := m["dest"][0]

	//accessing Google map api
	googleKey := os.Getenv("MAPS_KEY")
	gmap, err := maps.NewClient(maps.WithAPIKey(googleKey))

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DirectionsRequest{
		Origin:      origin,
		Destination: destination,
	}
	r.Mode = maps.TravelModeDriving
	r.Units = maps.UnitsMetric
	route, _, err := gmap.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	var maplegs *maps.Leg
	maplegs = route[0].Legs[0]
	pretty.Println(maplegs)

	var fullRoute models.Route
	db.DB.Create(&fullRoute)

	for i := 0; i < len(maplegs.Steps); i++ {

		polyline := maplegs.Steps[i].Polyline.Points
		var GooglePath models.GooglePath
		db.DB.Create(&GooglePath)

		GooglePath.RouteID = fullRoute.ID
		GooglePath.Path = polyline
		fullRoute.GooglePath = append(fullRoute.GooglePath, GooglePath)
		db.DB.Save(&GooglePath)
	}

	fullRoute.Duration = int(route[0].Legs[0].Duration.Minutes())
	fullRoute.Distance = route[0].Legs[0].Distance.Meters
	fullRoute.Fare = getFare(route[0].Legs[0].Distance.Meters)
	fullRoute.Source = origin
	fullRoute.Destination = destination
	db.DB.Save(&fullRoute)

	c.JSON(200, fullRoute)
}

func getFare(distance int) int {
	basePricePerMeter := 3
	return (distance / 1000) * basePricePerMeter

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

func testFromDB(c *gin.Context) {
	db := getDB(c)
	var taxis []models.Taxi
	db.Find(&taxis)

	fmt.Println(taxis)
	c.JSON(200, taxis)
}

func createRider(c *gin.Context) {
	var rider models.Rider
	db.DB.Create(&rider)

	c.JSON(200, rider.ID)
}

func bookingDBTest(c *gin.Context) {
	db := getDB(c)
	var bookings []models.Booking
	db.Find(&bookings)

	fmt.Println(bookings)
	c.JSON(200, bookings)
}
