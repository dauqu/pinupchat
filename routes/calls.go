package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pinupchat/actions"
	"pinupchat/config"
	"pinupchat/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CallsCollection *mongo.Collection = config.GetCollection(config.DB, "calls")

func AddCalls(c *gin.Context) {

	//Check Header Authorization
	if c.GetHeader("Authorization") == "" {
		c.JSON(400, gin.H{"message": "Invalid token"})
		return
	}

	idStr, err := actions.IdFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Read body
	var body models.Calls
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	from_id, _ := primitive.ObjectIDFromHex(idStr)
	to_id, _ := primitive.ObjectIDFromHex(body.To)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Add call to database
	_, err = CallsCollection.InsertOne(ctx, bson.M{
		"_id":        primitive.NewObjectID(),
		"from":       from_id,
		"to":         to_id,
		"type":       body.Type,
		"status":     body.Status,
		"is_deleted": false,
		"created_at": time.Now(),
		"updated_at": time.Now(),
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Room struct {
	ID           string
	Clients      map[*websocket.Conn]bool
	ClientsMutex sync.Mutex
}

var rooms map[string]*Room
var roomsMutex sync.Mutex

func createRoom(roomID string) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room := &Room{
		ID:           roomID,
		Clients:      make(map[*websocket.Conn]bool),
		ClientsMutex: sync.Mutex{},
	}
	rooms[roomID] = room
}

func deleteRoom(roomID string) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	delete(rooms, roomID)
}

func joinRoom(roomID string, conn *websocket.Conn) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room, ok := rooms[roomID]
	if !ok {
		log.Printf("Room '%s' not found\n", roomID)
		return
	}

	room.ClientsMutex.Lock()
	defer room.ClientsMutex.Unlock()

	room.Clients[conn] = true
	log.Printf("Client joined room '%s'\n", roomID)
}

func leaveRoom(roomID string, conn *websocket.Conn) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room, ok := rooms[roomID]
	if !ok {
		log.Printf("Room '%s' not found\n", roomID)
		return
	}

	room.ClientsMutex.Lock()
	defer room.ClientsMutex.Unlock()

	delete(room.Clients, conn)
	log.Printf("Client left room '%s'\n", roomID)
}



func CreateRooms(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Read room ID and join the room
	_, roomID, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	joinRoom(string(roomID), conn)

	// Read and write messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Broadcast message to all clients in the room
		room, ok := rooms[string(roomID)]
		if !ok {
			log.Printf("Room '%s' not found\n", roomID)
			break
		}

		room.ClientsMutex.Lock()
		for client := range room.Clients {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
				continue
			}
		}
		room.ClientsMutex.Unlock()
	}

	leaveRoom(string(roomID), conn)
}
