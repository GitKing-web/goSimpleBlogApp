package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username" validate:"required,min=3"`
	Email    string             `json:"email" bson:"email" validate:"required,min=5"`
	Password string             `json:"password" bson:"password" validate:"required"`
}
