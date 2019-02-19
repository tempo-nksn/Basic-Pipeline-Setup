package server

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"

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
	//db.DB.Create(&booking)

	booking.RiderID = uint(riderId)
	booking.TaxiID = uint(taxiId)
	booking.RouteID = uint(routeId)
	// 3) Get the duration, based on Routeid
	var route models.Route
	dbc := getDB(c)
	dbc.Where("id=?", routeId).Find(&route)
	booking.TravelDuration = route.Duration

	db.DB.Create(&booking)
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

func getDistance(c *gin.Context) {

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

	c.JSON(200, route[0].Legs[0].Distance.Meters)

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

func distPointLine(x float64, y float64, latLng maps.LatLng ) (int) {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * x / 180)
	radlat2 := float64(PI * latLng.Lat / 180)

	theta := float64(y - latLng.Lng)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	distInMeter := int(dist*1000)

	return distInMeter
}

func getRide(c *gin.Context) {
	database := getDB(c)

	r := c.Request
	m, _ := url.ParseQuery(r.URL.RawQuery)

	routeId,_ := strconv.Atoi(m["routeid"][0])
	var userRoute models.Route
	database.Where("id = ?", routeId).Preload("GooglePath").Find(&userRoute)

	userSrc := strings.Split(userRoute.Source, ",")
	userDest := strings.Split(userRoute.Destination, ",")

	var userSrcLat float64
	var userSrcLng float64
	var userDestLat float64
	var userDestLng float64

	userSrcLat, _ = strconv.ParseFloat(userSrc[0], 64)
	userSrcLng, _ = strconv.ParseFloat(userSrc[1], 64)
	userDestLat, _ = strconv.ParseFloat(userDest[0], 64)
	userDestLng, _ = strconv.ParseFloat(userDest[1], 64)

	var userRouteLatLang []maps.LatLng

	for _, googlePath := range userRoute.GooglePath {
		var latLang []maps.LatLng
		latLang, _ =  maps.DecodePolyline(googlePath.Path)
		userRouteLatLang = append(userRouteLatLang, latLang...)
	}

	var taxis []models.Taxi
	var route models.Route

	const statusActive  = "Active"
	const statusFree = "Free"
	const maxDistance  = 10 // Maximum distance of the user from a existing path for him/her to share a ride on that path

	var isSrcInRoute bool
	database.Where("status = ?", statusActive).Find(&taxis)

	if len(taxis) > 0 {
		// FInd if any existing active taxi can be shared
		for _, taxi := range taxis {

			database.Where("taxi_id = ? and status = ?", taxi.ID, statusActive).Preload("GooglePath").Find(&route)

			if route.ID == 0 {
				continue
			}

			routeDest := strings.Split(route.Destination, ",")
			var routeDestLat float64
			var routeDestLng float64

			routeDestLat, _ = strconv.ParseFloat(routeDest[0], 64)
			routeDestLng, _ = strconv.ParseFloat(routeDest[1], 64)

			var allLatLang []maps.LatLng

			for _, googlePath := range route.GooglePath {
				var latLang []maps.LatLng
				latLang, _ = maps.DecodePolyline(googlePath.Path)
				allLatLang = append(allLatLang, latLang...)
			}

			// Check iof the userSrc lies near the route
			isSrcInRoute = false
			for idx := 0; idx < len(allLatLang); idx++ {
				if distPointLine(userSrcLat, userSrcLng, allLatLang[idx]) < maxDistance {
					isSrcInRoute = true
					break
				}
			}

			// if the source is not near the route then check the next Taxi
			if !isSrcInRoute {
				continue
			}

			// If the source matches then check whether on eof the destination lies near the route of the other
			var isUserDestInRoute bool
			for idx := 0; idx < len(allLatLang); idx++ {
				if distPointLine(userDestLat, userDestLng, allLatLang[idx]) < maxDistance {
					isUserDestInRoute = true
					break
				}
			}

			if isUserDestInRoute {
				// Return the taxi id
				c.JSON(200, taxi.ID)
				return
			}

			// check if the taxi destination lies in userRoute
			var isRouteDestInUserRoute bool
			for idx := 0; idx < len(userRouteLatLang); idx++ {
				if distPointLine(routeDestLat, routeDestLng, userRouteLatLang[idx]) < maxDistance {
					isRouteDestInUserRoute = true
					break
				}
			}

			if isRouteDestInUserRoute {
				// Return the taxi id
				c.JSON(200, taxi.ID)
				return
			}
		}
	}

	// Give user the free taxi
	var taxi models.Taxi
	database.Where("status = ?", statusFree).First(&taxi)
	if taxi.NumberPlate != "" {
		c.JSON(200, taxi.ID)
		return
	}

	// No taxi Available
	c.JSON(400, "No Taxi available in the system")
}

func bookingDBTest(c *gin.Context) {
	db := getDB(c)
	var bookings []models.Booking
	db.Find(&bookings)

	fmt.Println(bookings)
	c.JSON(200, bookings)
}

func routeDBTest(c *gin.Context) {
	db := getDB(c)
	var route []models.Route
	db.Preload("GooglePath").Find(&route)

	fmt.Println(route)
	c.JSON(200, route)
}
