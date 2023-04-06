package main

import (
	"github.com/gin-gonic/gin"
	"pinupchat/routes"
)

func main() {
	//Create gin engine
	r := gin.Default()

	//Test route
	r.GET("/test", routes.InsertUser)

	//Group routes
	// auth := r.Group("/auth")
	// {
	// 	auth.POST("/login", Login)
	// 	auth.POST("/register", Register)
	// 	auth.GET("/profile", Profile)
	// }
	
	//Create http server
	r.Run(":8080")
}
