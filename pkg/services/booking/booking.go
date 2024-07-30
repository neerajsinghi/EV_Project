package booking

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/booking/db"
	"bikeRental/pkg/services/chronjobs"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"io"
	"net/http"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var service = db.NewService()

func Book(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w,r){
		return 
	}
	booking, err := getBooking(r)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.AddBooking(booking)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": err.Error()})
		return
	}
	chronjobs.CheckBooking()
	utils.SendOutput(err, w, r, data, booking, "Book")
}

func GetAllBookings(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w,r){
		return 
	}
	status := r.URL.Query().Get("status")
	bType := r.URL.Query().Get("type")
	vType := r.URL.Query().Get("vehicleType")
	data, err := service.GetAllBookings(status, bType, vType)
	utils.SendOutput(err, w, r, data, nil, "GetAllBookings")
}
func GetMyAllBooking(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w,r){
		return 
	}
	userID := mux.Vars(r)["id"]
	bType := r.URL.Query().Get("status")
	data, err := service.GetAllMyBooking(userID, bType)
	utils.SendOutput(err, w, r, data, nil, "GetMyAllBooking")
}
func GetMyLatestBooking(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w,r){
		return 
	}
	userID := mux.Vars(r)["id"]

	data, err := service.GetMyLatestBooking(userID)
	utils.SendOutput(err, w, r, data, nil, "GetMyLatestBooking")
}
func UpdateBooking(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w,r){
		return 
	}
	id := mux.Vars(r)["id"]
	booking, err := getBooking(r)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.UpdateBooking(id, booking)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to book"})
		return
	}
	dataBooking, _ := db.GetBooking(id)
	if booking.Status == "resumed" {
		dataBooking, _ = db.GetBooking(data)
	}
	utils.SendOutput(err, w, r, dataBooking, booking, "UpdateBooking")
}
func ResumeStoppedBooking(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w,r){
		return 
	}
	id := mux.Vars(r)["id"]
	data, err := service.UpdateBooking(id, entity.BookingDB{Status: "resumed"})
	if err != nil {
		trestCommon.ECLog1(err)
		return
	}
	chronjobs.CheckBooking()
	dataBooking, err := db.GetBooking(data)
	utils.SendOutput(err, w, r, dataBooking, "resuming", "ResumeStoppedBooking")
}
func GetBookingByID(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w,r){
		return 
	}
	id := mux.Vars(r)["id"]
	data, err := db.GetBooking(id)
	utils.SendOutput(err, w, r, data, nil, "GetBookingByID")
}
func getBooking(r *http.Request) (entity.BookingDB, error) {
	var user entity.BookingDB

	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}
	return user, err
}
