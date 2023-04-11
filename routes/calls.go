package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"pinupchat/actions"
	"pinupchat/config"
	"pinupchat/models"
	"time"
)

var CallsCollection *mongo.Collection = config.GetCollection(config.DB, "calls")

func AddCalls(c *gin.Context) {

	//Check Header Authorization
	if c.GetHeader("Authorization") == "" {
		c.JSON(400, gin.H{"message": "Invalid token"})
		return
	}

	idStr, err := actions.IdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Read body
	var body models.Calls
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	from_id, _ := primitive.ObjectIDFromHex(idStr)
	to_id, _ := primitive.ObjectIDFromHex(body.To)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Add call to database
	_, err = CallsCollection.InsertOne(ctx, bson.M{
		"_id":        primitive.NewObjectID(),
		"from":       from_id,
		"to":         to_id,
		"type":       body.Type,
		"status":     body.Status,
		"is_deleted": false,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Get status if deleted is false
	c.JSON(200, gin.H{"message": "Get calls"})
}

func GetCalls(c *gin.Context) {

	idStr, err := actions.IdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(idStr)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Get ID from param
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Get status if deleted is false
	var calls []bson.M
	err = CallsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&calls)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Get status if deleted is false
	c.JSON(200, gin.H{"message": "Get status"})
}
