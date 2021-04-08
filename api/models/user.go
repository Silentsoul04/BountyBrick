package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	Token     string             `json:"token"`
	LastLogin time.Time          `json:"last"`
	Role      string             `json:"role"`
}
