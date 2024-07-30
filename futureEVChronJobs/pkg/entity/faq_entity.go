package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type FAQDB struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Question    string             `json:"question" bson:"question"`
	Answer      string             `json:"answer" bson:"answer"`
	CreatedTime primitive.DateTime `json:"createdTime" bson:"created_time"`
}
