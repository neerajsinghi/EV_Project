package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookingDB represents the model for an booking
type BookingDB struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	ProfileID           string             `bson:"profile_id" json:"profileId"`
	DeviceID            int                `bson:"device_id" json:"deviceId"`
	StartTime           int64              `bson:"start_time" json:"startTime"`
	EndTime             int64              `bson:"end_time" json:"endTime"`
	StartKM             float64            `bson:"start_km" json:"startKm"`
	EndKM               float64            `bson:"end_km" json:"endKm"`
	TotalDistance       float64            `bson:"total_distance" json:"totalDistance"`
	Return              *Return            `bson:"return" json:"return"`
	Price               *float64           `bson:"price" json:"price"`
	Status              string             `bson:"status" json:"status"`
	VehicleType         string             `bson:"vehicle_type" json:"vehicleType"`
	BookingType         string             `bson:"booking_type" json:"bookingType"`
	Plan                *PlanDB            `bson:"plan" json:"plan"`
	CreatedTime         primitive.DateTime `bson:"created_time" json:"createdTime"`
	StartingStationID   string             `bson:"starting_station_id" json:"startingStationId"`
	EndingStationID     string             `bson:"ending_station_id" json:"endingStationId"`
	CarbonEmissionSaved float64            `bson:"carbon_emission_saved" json:"carbonEmissionSaved"`
	StartingStation     *StationDB         `bson:"starting_station" json:"startingStation"`
	EndingStation       *StationDB         `bson:"ending_station" json:"endingStation"`
	CouponCode          string             `bson:"coupon_code" json:"couponCode"`
	Discount            float64            `bson:"discount" json:"discount"`
	GreenPoints         int64              `bson:"green_points" json:"greenPoints"`
	CarbonSaved         float64            `bson:"carbon_saved" json:"carbonSaved"`
	City                string             `bson:"city" json:"city"`
}
type Return struct {
	Location      string   `bson:"location" json:"location"`
	Time          string   `bson:"time" json:"time"`
	ProductImages []string `bson:"product_images" json:"productImages"`
	Damages       []string `bson:"damages" json:"damages"`
	FrontPic      string   `bson:"front_picture" json:"frontPic"`
	BackPic       string   `bson:"back_picture" json:"backPic"`
	LeftPic       string   `bson:"left_picture" json:"leftPic"`
	RightPic      string   `bson:"right_picture" json:"rightPic"`
	FrontDesc     string   `bson:"front_desc" json:"frontDesc"`
	BackDesc      string   `bson:"back_desc" json:"backDesc"`
	LeftDesc      string   `bson:"left_desc" json:"leftDesc"`
	RightDesc     string   `bson:"right_desc" json:"rightDesc"`
}

type BookingOut struct {
	BookingDB      `json:",inline"`
	Profile        *ProfileDB `json:"profile,omitempty"`
	BikeWithDevice *IotBikeDB `json:"bikeWithDevice,omitempty"`
}
