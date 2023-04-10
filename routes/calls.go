package routes

import (
	"context"
	"fmt"
	"pinupchat/actions"
	"pinupchat/config"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CallsCollection *mongo.Collection = config.GetCollection(config.DB, "calls")

func AddCalls(c *gin.Context) {
	idStr, err := actions.IdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Add call to database
	_, err = CallsCollection.InsertOne(ctx, bson.M{
		"user_id":    idStr,
		"partner_id": c.Param("id"),
		"created_at": time.Now(),
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
