package faqdb

import (
	"bikeRental/pkg/entity"
	faqdb "bikeRental/pkg/services/faq/faqDB"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var service = faqdb.NewService()

func AddFaq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	faq, err := parseFaq(r)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{
			"status": false,
			"error":  "Something Went wrong",
		})
		return
	}
	data, err := service.AddFaq(faq)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{
			"status": false,
			"error":  "Unable to add FAQ",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{
		"status": true,
		"error":  "",
		"data":   data,
	})
}

func UpdateFaq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	faq, err := parseFaq(r)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{
			"status": false,
			"error":  "Something Went wrong",
		})
		return
	}
	id := mux.Vars(r)["id"]
	data, err := service.UpdateFaq(id, faq)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{
			"status": false,
			"error":  "Unable to update FAQ",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{
		"status": true,
		"error":  "",
		"data":   data,
	})
}
func DeleteFaq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := mux.Vars(r)["id"]
	err := service.DeleteFaq(id)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{
			"status": false,
			"error":  "Unable to delete FAQ",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{
		"status": true,
		"error":  "",
	})
}

func GetAllFaq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, err := service.GetAllFaq()
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{
			"status": false,
			"error":  "Unable to get FAQ",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{
		"status": true,
		"error":  "",
		"data":   data,
	})
}
func parseFaq(r *http.Request) (entity.FAQDB, error) {
	var faq entity.FAQDB
	err := json.NewDecoder(r.Body).Decode(&faq)
	if err != nil {
		return entity.FAQDB{}, err
	}
	return faq, nil
}
