package server

import "github.com/gin-gonic/gin"

func setupRoutes(router *gin.Engine) {
	authMiddleware := JWT()
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/signup", UserRegistration)

	v1 := router.Group("/api/v1")
	v1.GET("/", hello)
	v1.GET("/nearbytaxis", getNearByTaxis)
	v1.GET("/getRoute", getRoute)
	v1.GET("/userid", createRider)
	v1.GET("/ride", getRide)
	v1.POST("/bookingConfirmation", bookingConfirmation)
	v1.POST("/startRide", startRide)
	v1.POST("/finishRide", finishRide)
	v1.GET("/bookingDataBaseTest", bookingDBTest)
	v1.GET("/routeDBTest", routeDBTest)
	v1.GET("/getDistance", getDistance)
	v1.GET("/taxis", testFromDB)
	dash := v1.Group("/dashboard")
	dash.Use(authMiddleware.MiddlewareFunc())
	dash.GET("/", getUserDash)
}
