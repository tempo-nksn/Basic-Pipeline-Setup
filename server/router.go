package server

import "github.com/gin-gonic/gin"

func setupRoutes(router *gin.Engine) {
	//for local
	//path:= "templates/*"
	//for circleci
	path:="/go/src/github.com/tempo-nksn/Tempo-Backend/templates/*"
	router.LoadHTMLGlob(path)
	authMiddleware := JWT()
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/signup", userRegistration)
	router.GET("/templatetest", templateTest)
	driver:=router.Group("/driver")
	{
		driver.GET("/", driverIntro)
		driver.GET("/registration", driverReg)
		driver.POST("/registering", registering)
		driver.GET("/login", driverLogin)
		driver.POST("/dashboard", driverDash)
	}



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
	dash.Use(authMiddleware.LoginHandler)
	dash.GET("/", getUserDash)
}
