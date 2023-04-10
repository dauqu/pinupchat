package routes

import (
	"context"
	"net/http"
	"pinupchat/actions"
	"pinupchat/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ConversationCollection *mongo.Collection = config.GetCollection(config.DB, "conversation")

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func CreateMessage(c *gin.Context) {

	type Body struct {
		ConversationID string `json:"conversation_id"`
		Content        string `json:"content"`
	}

	//Create websocket connection
	ws, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		ws.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	// Parse request body
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		ws.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	// Check authorization token
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		ws.WriteJSON(gin.H{"error": "Invalid token"})
		return
	}

	userID, err := actions.IdFromToken(authToken)
	if err != nil {
		ws.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	// Convert conversation ID to ObjectID
	conversationID, err := primitive.ObjectIDFromHex(body.ConversationID)
	if err != nil {
		ws.WriteJSON(gin.H{"error": err.Error()})
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
		ws.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	// Send message to other user
	ws.WriteJSON(message)
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
