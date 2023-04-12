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

	poartnerid, _ := primitive.ObjectIDFromHex(body.PartnerId)
	userid, _ := primitive.ObjectIDFromHex(user_id)

	//CHeck partner id  and user id are not the same
	if user_id == body.PartnerId {
		c.JSON(400, gin.H{"message": "You can't add yourself"})
		return
	}

	//Check if partner exists
	_, err = UsersCollection.FindOne(context.Background(), bson.M{"_id": poartnerid}).DecodeBytes()
	if err != nil {
		c.JSON(400, gin.H{"message": "Partner not found"})
		return
	}

	//Check if conversation already exists
	_, err = ConversationCollection.FindOne(context.Background(), bson.M{
		"participents": bson.A{
			bson.M{
				"_id": userid,
			},
			bson.M{
				"_id": poartnerid,
			},
		},
	}).DecodeBytes()
	if err == nil {
		c.JSON(400, gin.H{"message": "Conversation already exists"})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Insert conversation
	_, err = ConversationCollection.InsertOne(ctx, bson.M{
		"participents": bson.A{
			bson.M{
				"_id": userid,
			},
			bson.M{
				"_id": poartnerid,
			},
		},
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
	// Create the $or operator
	filter := bson.M{
		"participents._id": userid,
		// "$or": bson.A{
		// 	bson.M{"partner_id": userid},
		// 	bson.M{"user_id": userid},
		// },
	}

	//pipeline
	pipeline := bson.A{
		bson.M{"$match": filter},
		// bson.M{"$lookup": bson.M{
		// 	"from":         "users",
		// 	"localField":   "partner_id",
		// 	"foreignField": "_id",
		// 	"as":           "partner",
		// }},
		// bson.M{"$lookup": bson.M{
		// 	"from":         "users",
		// 	"localField":   "user_id",
		// 	"foreignField": "_id",
		// 	"as":           "user",
		// }},
		//Lookup participents and get their names
		bson.M{"$lookup": bson.M{
			"from":         "users",
			"localField":   "participents._id",
			"foreignField": "_id",
			"as":           "participents",
		}},
		//Remove password from user object
		bson.M{"$project": bson.M{
			"participents.password": 0,
		}},

		//Remove messages from conversation
		bson.M{"$project": bson.M{
			"messages": 0,
		}},

		//Remove blocked conversations
		bson.M{"$match": bson.M{
			"blocked": false,
		}},

		//Remove archived conversations
		bson.M{"$match": bson.M{
			"archived": false,
		}},

		//Remove my ID from participents
		bson.M{"$project": bson.M{
			"participents": bson.M{
				"$filter": bson.M{
					"input": "$participents",
					"as":    "participent",
					"cond": bson.M{
						"$ne": bson.A{"$$participent._id", userid},
					},
				},
			},
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

	c.JSON(200, gin.H{"user_id": user_id, "conversations": conversations})
}

func Archived(c *gin.Context) {

}

func Delete(c *gin.Context) {

}

func Block(c *gin.Context) {

}

func Mute(c *gin.Context) {

}
