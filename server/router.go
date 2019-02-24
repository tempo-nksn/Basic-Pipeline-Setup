package server

import "github.com/gin-gonic/gin"

func setupRoutes(router *gin.Engine) {
	router.LoadHTMLGlob("../templates/*")
	authMiddleware := JWT()
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/signup", UserRegistration)
	router.GET("/templatetest", TemplateTest)
	driver:=router.Group("/driver")
	{
		driver.GET("/", Driver)
		driver.GET("/registration", DriverReg)
		driver.POST("/registering", Registering)
		driver.GET("/login", DriverLogin)
		driver.POST("/dashboard", DriverDash)

	}



	v1 := router.Group("/api/v1")
	v1.GET("/", Hello)
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
	dash.Use(authMiddleware.LoginHandler)
	dash.GET("/", getUserDash)
}
