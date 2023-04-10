package routes

import (
	"context"
	"pinupchat/actions"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	userid, _ := primitive.ObjectIDFromHex(user_id)
	oartnerid, _ := primitive.ObjectIDFromHex(body.PartnerId)

	//Insert conversation
	_, err = ConversationCollection.InsertOne(ctx, bson.M{
		"user_id":    userid,
		"partner_id": oartnerid,
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

	userid, _ := primitive.ObjectIDFromHex(user_id)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Filter
	filter := bson.M{
		"$or": []bson.M{
			{"user_id": userid},
			{"partner_id": userid},
		},
	}

	//pipeline
	pipeline := bson.A{
		bson.M{"$match": filter},
		bson.M{"$lookup": bson.M{
			"from":         "users",
			"localField":   "partner_id",
			"foreignField": "_id",
			"as":           "partner",
		}},
		// bson.M{"$lookup": bson.M{
		// 	"from":         "users",
		// 	"localField":   "user_id",
		// 	"foreignField": "_id",
		// 	"as":           "user",
		// }},

		//Remove password from user object
		bson.M{"$project": bson.M{
			"partner.password": 0,
		}},
		//Remove messages from conversation
		bson.M{"$project": bson.M{
			"messages": 0,
		}},
	}

	//Get conversation
	cursor, err := ConversationCollection.Aggregate(ctx, pipeline)
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
