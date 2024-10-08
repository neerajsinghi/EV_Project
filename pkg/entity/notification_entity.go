package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	Title       string             `json:"title" bson:"title"`
	Body        string             `json:"body" bson:"body"`
	UserId      string             `json:"userId" bson:"userId"`
	Token       string             `json:"token" bson:"token"`
	Type        string             `json:"type" bson:"type"`
	CreatedTime primitive.DateTime `json:"createdTime" bson:"created_time"`
}

type NotificationMulti struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Body        string             `json:"body" bson:"body"`
	UserIds     []string           `json:"userId" bson:"userIds"`
	Tokens      []string           `json:"token" bson:"tokens"`
	Type        string             `json:"type" bson:"type"`
	CreatedTime primitive.DateTime `json:"createdTime" bson:"created_time"`
}

type PreDefNotification struct {
	Name  string `json:"name" bson:"name"`
	Title string `json:"title" bson:"title"`
	Body  string `json:"body" bson:"body"`
	Type  string `json:"type" bson:"type"`
}
