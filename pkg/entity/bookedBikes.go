package entity

type BookedBikesDB struct {
	UserID    string    `json:"userId" bson:"user_id"`
	BookingID string    `json:"bookingId" bson:"booking_id"`
	Bike      IotBikeDB `json:"bike" bson:"bike"`
	OnGoing   bool      `json:"onGoing" bson:"on_going"`
	Profile   ProfileDB `json:"profile" bson:"profile"`
	Booking   BookingDB `json:"booking" bson:"booking"`
}
