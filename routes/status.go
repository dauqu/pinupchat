package routes

import (
	"context"
	"pinupchat/config"
	"pinupchat/models"
	"time"

	"pinupchat/actions"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var StatusCollection *mongo.Collection = config.GetCollection(config.DB, "status")

func CreateStatus(c *gin.Context) {

	id, err := actions.IdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var status models.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Insert status
	_, err = StatusCollection.InsertOne(ctx, bson.M{
		"user_id":    id,
		"content":    status.Content,
		"is_deleted": false,
		"seen_by":    bson.A{},
		"created_at": time.Now(),
	})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Status created"})
}

func GetStatus(c *gin.Context) {
	idStr, err := actions.IdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Get status if deleted is false
	cursor, err := StatusCollection.Find(ctx, bson.M{"user_id": idStr, "is_deleted": false})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Status found", "status": results})
}

// Delete status
func DeleteStatus(c *gin.Context) {

	user_id, err := actions.IdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	statusId := c.Param("id")

	id, _ := primitive.ObjectIDFromHex(statusId)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Delete status
	_, err = StatusCollection.UpdateOne(ctx, bson.M{"_id": id, "user_id": user_id}, bson.M{"$set": bson.M{"is_deleted": true}})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Status deleted"})
}
