package entity

type BookedBikesDB struct {
	UserID    string     `json:"userId" bson:"user_id"`
	BookingID string     `json:"bookingId" bson:"booking_id"`
	OnGoing   bool       `json:"onGoing" bson:"on_going"`
	Booking   BookingOut `json:"booking" bson:"booking"`
}
