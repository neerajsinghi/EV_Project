package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type ServiceDB struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Type        serviceType        `json:"type" bson:"type"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Active      bool               `json:"active" bson:"active"`
	Discount    float64            `json:"discount" bson:"discount"`
	Status      string             `json:"status" bson:"status"`
	CreatedTime primitive.DateTime `json:"createdTime" bson:"created_time"`
}

type serviceType string

const (
	Hourly   serviceType = "hourly"
	Rental   serviceType = "rental"
	Charging serviceType = "charging"
	ECar     serviceType = "eCar"
)
