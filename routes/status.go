package routes

import (
	"pinupchat/config"

	"go.mongodb.org/mongo-driver/mongo"
)

var StatusCollection *mongo.Collection = config.GetCollection(config.DB, "status")

func Status() {
	//Get my status from the database	

	//Get my status from the database
}