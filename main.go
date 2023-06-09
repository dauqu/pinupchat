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

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to PinupChat",
		})
	})

	//Group routes
	autha := r.Group("/auth")
	{
		autha.POST("/register", auth.Register)
		autha.POST("/login", auth.Login)
		autha.GET("/check-login", auth.CheckLogin)
		autha.GET("/profile", auth.Profile)
		autha.GET("/get-all-users", auth.GetAllUsers)
	}

	routess := r.Group("/api")
	{
		routess.POST("/add-status", routes.CreateStatus)
		routess.GET("/get-status", routes.GetStatus)
		routess.DELETE("/delete-status/:id", routes.DeleteStatus)

		routess.POST("/create-contact", routes.CreateContact)
		routess.GET("/get-contacts", routes.GetContacts)

		routess.Any("/create-message", routes.CreateMessage)
		routess.POST("/delete-message", routes.DeleteMessage)
		routess.GET("/get-messages/:id", routes.GetMessages)

		routess.POST("/add-calls", routes.AddCalls)
		routess.GET("/get-calls", routes.GetCalls)

		routess.Any("/create-rooms", routes.CreateRooms)

	}

	//Create http server
	r.Run(":8080")
}
