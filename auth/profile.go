package auth

import (
	"context"
	"fmt"
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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Check if token is valid
	if token.Valid {
		fmt.Println(claims["id"])
	}

	//MongoDb primitive object id
	id, _ := primitive.ObjectIDFromHex(claims["id"].(string))

	cursor, err := UsersCollection.Find(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User found", "user": results})
}
