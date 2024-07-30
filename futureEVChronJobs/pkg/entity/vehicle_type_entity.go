package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type VehicleTypeDB struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Price       float64            `bson:"price" json:"price"`
	Description string             `bson:"description" json:"description"`
	CreatedTime primitive.DateTime `bson:"created_time" json:"created_time"`
}
