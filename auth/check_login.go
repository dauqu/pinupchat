package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckLogin(c *gin.Context) {

	// Get authorization header
	tokenString := c.Request.Header.Get("Authorization")

	// Check if token is empty
	if tokenString == "" || tokenString == "null" || tokenString == "undefined" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token is empty", "isLogged": false})
		return
	}

	// Verify token and return claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// TODO: Replace with your own secret key
		return []byte("secret"), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "isLogged": false})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token is invalid", "isLogged": false})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token claims", "isLogged": false})
		return
	}

	idStr, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID", "isLogged": false})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID", "isLogged": false})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user bson.D
	err = UsersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found", "isLogged": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error", "isLogged": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token is valid", "isLogged": true})
}
