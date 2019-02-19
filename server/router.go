package server

import "github.com/gin-gonic/gin"

func setupRoutes(router *gin.Engine) {

	v1 := router.Group("/api/v1")

	v1.GET("/", hello)
	v1.POST("/", dummyPost)
	v1.GET("/nearbytaxis/", getNearByTaxis)
	v1.GET("/getRoute", getRoute)
	v1.GET("/taxis/", testFromDB)
	v1.GET("/userid", createRider)

	v1.POST("/bookingConfirmation", bookingConfirmation)
	v1.GET("/bookingDataBaseTest", bookingDBTest)
	v1.GET("/routeDBTest"routeDBTest)

	//v1.GET("/ride/", getRide)

}
