package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID       primitive.ObjectID `bson:"_id"`
	Data     string             `bson:"data"`
	Created  time.Time          `bson:"created"`
	Modified time.Time          `bson:"modified"`
}
