package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Posts struct {
	UserId    primitive.ObjectID  `json:"userId" bson:"userId"`
	PostId    primitive.ObjectID  `json:"postId" bson:"postId"`
	Title     string              `json:"title" bson:"title"`
	Content   string              `json:"content" bson:"content"`
	CreatedAt primitive.Timestamp `json:"createdAt" bson:"createdAt"`
}
