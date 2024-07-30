package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Stations maps to bike

type StationDB struct {
	ID                primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name              string             `bson:"name" json:"name,omitempty"`
	Description       string             `bson:"description" json:"description,omitempty"`
	ShortName         string             `bson:"short_name" json:"shortName,omitempty"`
	Address           *AddressDB         `bson:"address" json:"address,omitempty"`
	Location          *Location          `bson:"location" json:"location,omitempty"`
	Active            *bool              `bson:"active" json:"active,omitempty"`
	Group             string             `bson:"group" json:"group,omitempty"`
	SupervisorID      string             `bson:"supervisor_id" json:"supervisorID,omitempty"`
	Stock             *int               `bson:"stock" json:"stock,omitempty"`
	Public            *bool              `bson:"public" json:"public,omitempty"`
	Status            string             `bson:"status" json:"status,omitempty"`
	ServicesAvailable []string           `bson:"services_available" json:"servicesAvailable,omitempty"`
	UpdateAt          primitive.DateTime `bson:"update_at" json:"updateAt,omitempty"`
	CreatedTime       primitive.DateTime `bson:"created_time" json:"createdTime,omitempty"`
	LocationPolygon   *LocationPolygon   `bson:"location_polygon" json:"locationPolygon,omitempty"`
}

// City, Type, Station, Vehicle Type, IOT device data, upload insurance policy, insurance date which gives notification to the admin once its time to renew, vehicle registration, any permits required to deploy the vehicle etc. (by clicking on the vehicle you can check all the details about it, number of ride, ride history, money earned, current location, Vehicle status (as shared earlier in the group) (1.⁠ ⁠On Trip - Any vehicle that is currently booked  2.⁠ ⁠Available - Any vehicle that is active and available for the user to be booked  3.⁠ ⁠Damaged - Any vehicle that has been marked damaged by the ops team 4.⁠ ⁠Under Maintenance - Any vehicle that has been gone for servicing or repair 5.⁠ Ready for Deployment - First status after we list the vehicle on the system  (provide an option on the admin panel to manually change the status of the vehicle) (These statuses can be changed by the ops, we should know who are marking these changes) (Maintenance chart needs to be maintained about each vehicle)
type DeviceInfo struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	City                string             `bson:"city" json:"city"`
	Type                string             `bson:"type" json:"type"`
	DeviceID            int                `bson:"device_id" json:"deviceId,omitempty"`
	VehicleTypeID       string             `bson:"vehicle_type_id" json:"vehicleTypeId,omitempty"`
	StationID           *string            `bson:"station_id" json:"stationId,omitempty"`
	Status              string             `bson:"status" json:"status,omitempty"`
	VehicleType         *VehicleTypeDB     `bson:"vehicle_type" json:"vehicleType,omitempty"`
	DeviceData          *IotBikeDB         `bson:"device_data" json:"deviceData,omitempty"`
	Station             *StationDB         `bson:"station" json:"station,omitempty"`
	Stations            []StationDB        `bson:"stations" json:"stations,omitempty"`
	Description         string             `bson:"description" json:"description,omitempty"`
	CreatedTime         primitive.DateTime `bson:"created_time" json:"createdTime,omitempty"`
	InsuranceDate       time.Time          `bson:"insurance_date" json:"insuranceDate,omitempty"`
	InsurancePolicy     string             `bson:"insurance_policy" json:"insurancePolicy,omitempty"`
	VehicleRegistration string             `bson:"vehicle_registration" json:"vehicleRegistration,omitempty"`
	PermitsRequired     []string           `bson:"permits_required" json:"permitsRequired,omitempty"`
	Immobilized         bool               `bson:"immobilized" json:"immobilized,omitempty"`
}
