package booking

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/booking/db"
	"bikeRental/pkg/services/chronjobs"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"io"
	"net/http"
	"time"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var service = db.NewService()

func Book(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.WriteHeader(http.StatusOK)
	chronjobs.CheckBooking()
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}

func GetAllBookings(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	status := r.URL.Query().Get("status")
	bType := r.URL.Query().Get("type")
	vType := r.URL.Query().Get("vehicleType")
	data, err := service.GetAllBookings(status, bType, vType)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to find brand"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}
func GetMyAllBooking(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	userID := mux.Vars(r)["id"]
	bType := r.URL.Query().Get("status")
	data, err := service.GetAllMyBooking(userID, bType)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to find brand"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}
func GetMyLatestBooking(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	userID := mux.Vars(r)["id"]

	data, err := service.GetMyLatestBooking(userID)
	if err != nil || data == nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to find latest booking"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}
func UpdateBooking(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.WriteHeader(http.StatusOK)
	dataBooking, _ := db.GetBooking(id)
	if booking.Status == "resumed" {
		dataBooking, _ = db.GetBooking(data)
	}

	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": dataBooking, "dataSucces": data})
}
func ResumeStoppedBooking(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	data, err := service.UpdateBooking(id, entity.BookingDB{Status: "resumed"})
	if err != nil {
		trestCommon.ECLog1(err)
		return
	}
	chronjobs.CheckBooking()
	dataBooking, _ := db.GetBooking(data)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": dataBooking})
}
func GetBookingByID(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := mux.Vars(r)["id"]
	data, err := db.GetBooking(id)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": true, "error": "unable to find booking", "data": data})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": false, "error": "", "data": data})
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
