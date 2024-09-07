package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Posts struct {
	PostId  primitive.ObjectID
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
}
