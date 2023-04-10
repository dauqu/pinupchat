package routes

import (
	"context"
	"net/http"
	"pinupchat/actions"
	"pinupchat/config"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ConversationCollection *mongo.Collection = config.GetCollection(config.DB, "conversation")

func CreateMessage(c *gin.Context) {
	type Body struct {
		ConversationID string `json:"conversation_id"`
		Content        string `json:"content"`
	}

	// Parse request body
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check authorization token
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization token"})
		return
	}

	userID, err := actions.IdFromToken(authToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization token"})
		return
	}

	// Convert conversation ID to ObjectID
	conversationID, err := primitive.ObjectIDFromHex(body.ConversationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Create message object
	message := bson.M{
		"sender":     userID,
		"content":    body.Content,
		"is_read":    false,
		"is_deleted": false,
		"is_edited":  false,
		"created_at": time.Now(),
		"updated_at": "",
	}

	// Add message to conversation
	filter := bson.M{"_id": conversationID}
	update := bson.M{"$push": bson.M{"messages": message}}
	_, err = ConversationCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add message to conversation"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Message created"})
}

func GetMessages(c *gin.Context) {

	//Return if Auth token is invalid
	// if c.GetHeader("Authorization") == "" || c.GetHeader("Authorization") == "null" || c.GetHeader("Authorization") == "undefined" {
	// 	c.JSON(400, gin.H{"message": "Invalid token"})
	// 	return
	// }

	// user_id, err := actions.IdFromToken(c.GetHeader("Authorization"))
	// if err != nil {
	// 	c.JSON(400, gin.H{"message": err.Error()})
	// 	return
	// }

	//Get ID from params
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": id}

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

	c.JSON(200, conversations)
}
