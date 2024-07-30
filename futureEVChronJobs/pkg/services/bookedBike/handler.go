package bookedbike

import (
	bookedlogic "futureEVChronJobs/pkg/services/bookedBike/logic"
	utils "futureEVChronJobs/pkg/util"
	"net/http"
)

func GetBookedBike(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	userID := r.URL.Query().Get("userID")
	bookingId := r.URL.Query().Get("bookingId")
	data, err := bookedlogic.GetBookedBike(userID, bookingId)
	utils.SendOutput(err, w, r, data, nil, "GetBookedBike")
}
