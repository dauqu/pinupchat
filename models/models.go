package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Phone        string             `bson:"phone,omitempty"`
	Email        string             `bson:"email,omitempty"`
	Country      string             `bson:"country,omitempty"`
	About        string             `bson:"about,omitempty"`
	Avatar       string             `bson:"avatar,omitempty"`
	OnlineStatus string             `bson:"online_status,omitempty"` // online, offline, busy
	LastActive   time.Time          `bson:"last_active,omitempty"`   // last active time
	Password     string             `bson:"password,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
}

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Admin     string             `bson:"admin,omitempty"` // admin id
	Pic       string             `bson:"pic,omitempty"`   // url
	About     string             `bson:"about,omitempty"` // room description
	Name      string             `bson:"name,omitempty"`
	Users     []string           `bson:"users,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

type Status struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id,omitempty"` // user id
	Content   string             `bson:"content,omitempty"`
	IsDeleted bool               `bson:"is_deleted,omitempty"`
	SeenBy    []string           `bson:"seen_by,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
}

type Calls struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	From      string             `bson:"from,omitempty"` // user id
	To        string             `bson:"to,omitempty"`   // user id or room id
	IsVideo   bool               `bson:"is_video,omitempty"`
	IsDeleted bool               `bson:"is_deleted,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

type ContectList struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id,omitempty"` // user id
	Contects  []string           `bson:"contects,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

type Group struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Admin     string             `bson:"admin,omitempty"` // admin id
	Pic       string             `bson:"pic,omitempty"`   // url
	About     string             `bson:"about,omitempty"` // room description
	Name      string             `bson:"name,omitempty"`
	Users     []string           `bson:"users,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

type Conversation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"` // user id
	PartnerID primitive.ObjectID `bson:"partner_id,omitempty"`
	Archived  bool               `bson:"archived,omitempty"`
	Deleted   bool               `bson:"deleted,omitempty"`
	Blocked   bool               `bson:"blocked,omitempty"`
	Muted     bool               `bson:"muted,omitempty"`
	Messages  []Message          `bson:"messages,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Sender    string             `bson:"sender,omitempty"` // user id
	Content   string             `bson:"content,omitempty"`
	IsDeleted bool               `bson:"is_deleted,omitempty"`
	IsEdited  bool               `bson:"is_edited,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}
