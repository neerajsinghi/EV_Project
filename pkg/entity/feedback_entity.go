package entity

type Feedback struct {
	Feedback  string `json:"feedback" bson:"feedback"`
	Ratings   int    `json:"ratings" bson:"ratings"`
	ProfileID string `json:"profile_id" bson:"profile_id"`
	BookingID string `json:"booking_id" bson:"booking_id"`
}
type FeedbackOut struct {
	Feedback  string    `json:"feedback" bson:"feedback"`
	Ratings   int       `json:"ratings" bson:"ratings"`
	ProfileID string    `json:"profile_id" bson:"profile_id"`
	BookingID string    `json:"booking_id" bson:"booking_id"`
	Booking   BookingDB `json:"booking" bson:"booking"`
	Profile   ProfileDB `json:"profile" bson:"profile"`
}

type Response struct {
	Status bool        `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}
