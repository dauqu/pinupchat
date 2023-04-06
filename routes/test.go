package routes

import (
	"context"
	"pinupchat/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var UsersCollection *mongo.Collection = config.GetCollection(config.DB, "users")

func Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// Insert Users
func InsertUser(c *gin.Context) {

	resp, err := UsersCollection.InsertOne(context.Background(), bson.M{
		"username": "test",
		"password": "test",
	})
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error inserting user",
		})
	}
	c.JSON(200, gin.H{
		"message":  "User inserted",
		"response": resp,
	})
}
