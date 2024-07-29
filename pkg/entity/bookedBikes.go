package entity

type BookedBikesDB struct {
	UserID          string     `json:"userId" bson:"user_id"`
	BookingID       string     `json:"bookingId" bson:"booking_id"`
	OnGoing         bool       `json:"onGoing" bson:"on_going"`
	Booking         BookingOut `json:"booking" bson:"booking"`
	Coordinates     []float64  `json:"coordinates" bson:"coordinates"`
	UserName        string     `json:"userName" bson:"user_name"`
	StartingStation string     `json:"startingStation" bson:"starting_station"`
	EndStation      string     `json:"endStation" bson:"end_station"`
	DeviceName      string     `json:"deviceName" bson:"device_name"`
	DeviceId        int        `json:"deviceId" bson:"device_id"`
}
