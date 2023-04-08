package auth

import (
	"context"
	"pinupchat/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		//MongoDb primitive object id
		id, _ := primitive.ObjectIDFromHex(claims["id"].(string))

		cursor, err := UsersCollection.Find(ctx, bson.M{"_id": id})
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		var users []models.User
		if err = cursor.All(ctx, &users); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Token is valid", "isLogged": true, "user": users})
	}

	c.JSON(200, gin.H{"message": "Token is invalid", "isLogged": true})
}
