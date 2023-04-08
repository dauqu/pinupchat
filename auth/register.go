package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"pinupchat/models"
	"time"
)

func Register(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Check all fields
	if user.FullName != "" || user.Phone != "" || user.Password != "" || user.Country != "" {
		c.JSON(400, gin.H{"message": "All fields are required"})
		return
	}

	//Check if username exists
	var check_user models.User
	err := UsersCollection.FindOne(ctx, bson.M{"phone": user.Phone}).Decode(&check_user)
	if err != mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"message": "Username already exists"})
		return
	}

	// //Check if email exists
	// var check_email models.User
	// err = UsersCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&check_email)
	// if err != mongo.ErrNoDocuments {
	// 	c.JSON(400, gin.H{"message": "Email already exists"})
	// 	return
	// }

	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Insert user
	_, err = UsersCollection.InsertOne(ctx, bson.M{
		"full_name":     user.FullName,
		"phone":         user.Phone,
		"country":       user.Country,
		"email":         "",
		"about":         "",
		"avatar":        "",
		"online_status": "offline",
		"last_active":   time.Now(),
		"password":      hashedPassword,
		"created_at":    time.Now(),
	})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "User created",
	})
}
