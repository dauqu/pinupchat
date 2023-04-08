package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Profile(c *gin.Context) {

	//Get authorization header
	token := c.Request.Header.Get("Authorization")

	//Check if token is empty
	if token == "" || token == "null" || token == "undefined" {
		c.JSON(200, gin.H{"message": "Token is empty", "isLogged": false})
		return
	}

	//Verify token and return claims
	claims, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		c.JSON(200, gin.H{"message": err.Error(), "isLogged": false})
		return
	}

	//Check if token is valid
	if claims.Valid {
		//Get claims
		claims := claims.Claims.(*jwt.MapClaims)
		c.JSON(200, gin.H{"message": "Token is valid", "isLogged": true, "claims": claims})
	}

	c.JSON(200, gin.H{"message": "Token is invalid", "isLogged": true})
}