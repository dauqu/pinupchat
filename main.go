package main

import (
	// "fmt"
	// "pinupchat/actions"
	// "time"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v4"
	"pinupchat/auth"
)

func main() {

	//Create gin engine and routes
	r := gin.Default()

	//Generate token
	// token := jwt.New(jwt.SigningMethodHS256)
	// claims := token.Claims.(jwt.MapClaims)
	// claims["username"] = "Harsha Web"
	// claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// t, err := token.SignedString([]byte("secret"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	
	// fmt.Println(t)
	
	// ttt, err := actions.VerifyToken(t)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(ttt)

	//Group routes
	autha := r.Group("/auth")
	{
		autha.POST("/register", auth.Register)
		autha.POST("/login", auth.Login)
		autha.GET("/check-login", auth.CheckLogin)
		autha.GET("/profile", auth.Profile)
	}
	
	//Create http server
	r.Run(":8080")
}
