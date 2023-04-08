package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Profile(c *gin.Context) {
	// Get authorization header
	tokenString := c.Request.Header.Get("Authorization")

	// Check if token is empty
	if tokenString == "" || tokenString == "null" || tokenString == "undefined" {
		c.JSON(200, gin.H{"message": "Token is empty", "isLogged": false})
		return
	}

	// Parse and verify the token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		c.JSON(200, gin.H{"message": err.Error(), "isLogged": false})
		return
	}

	// Check if token is valid
	if token.Valid {
		c.JSON(200, gin.H{"message": "Token is valid", "isLogged": true, "claims": claims})
		return
	}

	c.JSON(200, gin.H{"message": "Token is invalid", "isLogged": true})
}
