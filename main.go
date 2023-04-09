package main

import (
	// "fmt"
	// "pinupchat/actions"
	// "time"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v4"
	"pinupchat/auth"
	"pinupchat/routes"
)

func main() {

	//Create gin engine and routes
	r := gin.Default()

	//Group routes
	autha := r.Group("/auth")
	{
		autha.POST("/register", auth.Register)
		autha.POST("/login", auth.Login)
		autha.GET("/check-login", auth.CheckLogin)
		autha.GET("/profile", auth.Profile)
	}

	routess := r.Group("/api")
	{
		routess.POST("/add-status", routes.CreateStatus)
		routess.GET("/get-status", routes.GetStatus)
		routess.DELETE("/delete-status/:id", routes.DeleteStatus)

	}
	
	//Create http server
	r.Run(":8080")
}
