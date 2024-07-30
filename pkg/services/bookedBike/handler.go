package bookedbike

import (
	bookedlogic "bikeRental/pkg/services/bookedBike/logic"
	utils "bikeRental/pkg/util"
	"net/http"
)

func GetBookedBike(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	userID := r.URL.Query().Get("userID")
	bookingId := r.URL.Query().Get("bookingId")
	data, err := bookedlogic.GetBookedBike(userID, bookingId)
	utils.SendOutput(err, w, r, data, nil, "GetBookedBike")
}
