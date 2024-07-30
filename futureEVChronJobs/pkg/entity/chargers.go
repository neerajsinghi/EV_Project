package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChargerDB struct {
	ID           primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name,omitempty"`
	Description  string             `bson:"description" json:"description,omitempty"`
	ShortName    string             `bson:"short_name" json:"shortName,omitempty"`
	Address      *AddressDB         `bson:"address" json:"address,omitempty"`
	Location     *Location          `bson:"location" json:"location,omitempty"`
	Active       *bool              `bson:"active" json:"active,omitempty"`
	SupervisorID string             `bson:"supervisor_id" json:"supervisorID,omitempty"`
	Stock        *int               `bson:"stock" json:"stock,omitempty"`
	Public       *bool              `bson:"public" json:"public,omitempty"`
	Status       string             `bson:"status" json:"status,omitempty"`
	CreatedTime  primitive.DateTime `bson:"created_time" json:"createdTime,omitempty"`
}
