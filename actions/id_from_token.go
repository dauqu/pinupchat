package actions

import (
	// "fmt"
	"context"
	"pinupchat/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var UsersCollection *mongo.Collection = config.GetCollection(config.DB, "users")

func IdFromToken(tokenString string) (res string, err error) {

	// Check if token is empty
	if tokenString == "" || tokenString == "null" || tokenString == "undefined" {
		return "invalid token", nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// TODO: Replace with your own secret key
		return []byte("secret"), nil
	})
	if err != nil {
		return "invalid token", nil
	}

	if !token.Valid {
		return "invalid token", nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "invalid token", nil
	}

	idStr, ok := claims["id"].(string)
	if !ok {
		return "invalid token", nil
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return "invalid token", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user bson.D
	err = UsersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "user not found", nil
		}
		return err.Error(), err
	}

	return idStr, nil
}
