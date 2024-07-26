package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlanDB struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Name            string             `bson:"name" json:"name"`
	City            string             `bson:"city" json:"city"`
	VehicleType     string             `bson:"vehicle_type" json:"vehicleType"`
	ChargerType     string             `bson:"charger_type" json:"chargerType"`
	Type            serviceType        `json:"type" bson:"type"`
	Description     string             `bson:"description" json:"description"`
	StartingMinutes int                `bson:"starting_minutes" json:"startingMinutes"`
	EndingMinutes   int                `bson:"ending_minutes" json:"endingMinutes"`
	EveryXMinutes   int                `bson:"every_x_minutes" json:"everyXMinutes"`
	Price           float64            `bson:"price" json:"price"`
	Deposit         *float64           `bson:"deposit" json:"deposit"`
	Validity        string             `bson:"validity" json:"validity"`
	Discount        float64            `bson:"discount" json:"discount"`
	IsActive        *bool              `bson:"is_active" json:"isActive"`
	Status          string             `bson:"status" json:"status"`
	CreatedTime     primitive.DateTime `bson:"created_time" json:"createdTime"`
}
