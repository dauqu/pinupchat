package routes

import (
	"context"
	"pinupchat/actions"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateContact(c *gin.Context) {

	type Body struct {
		PartnerId string `json:"partner_id"`
	}
	//Bind JSON
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Return if Auth token is invalid
	if c.GetHeader("Authorization") == "" || c.GetHeader("Authorization") == "null" || c.GetHeader("Authorization") == "undefined" {
		c.JSON(400, gin.H{"message": "Invalid token"})
		return
	}

	user_id, err := actions.IdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Insert conversation
	_, err = ConversationCollection.InsertOne(ctx, bson.M{
		"user_id":    user_id,
		"partner_id": body.PartnerId,
		"messages":   bson.A{},
		"archived":   false,
		"deleted":    false,
		"blocked":    false,
		"muted":      false,
		"created_at": time.Now(),
	})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Conversation created"})
}

func GetContacts(c *gin.Context) {
	//Return if Auth token is invalid
	if c.GetHeader("Authorization") == "" || c.GetHeader("Authorization") == "null" || c.GetHeader("Authorization") == "undefined" {
		c.JSON(400, gin.H{"message": "Invalid token"})
		return
	}

	user_id, err := actions.IdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{
		"$or": []bson.M{
			{"user_id": user_id},
			{"partner_id": user_id},
		},
	}

	//Get conversation
	cursor, err := ConversationCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var conversations []bson.M
	if err = cursor.All(ctx, &conversations); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Conversations fetched", "id": user_id, "conversations": conversations})
}
