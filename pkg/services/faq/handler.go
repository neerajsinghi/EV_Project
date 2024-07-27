package faqdb

import (
	"bikeRental/pkg/entity"
	faqdb "bikeRental/pkg/services/faq/faqDB"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var service = faqdb.NewService()

func AddFaq(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
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
	utils.SendOutput(err, w, r, data, "AddFaq")
}

func UpdateFaq(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
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
	utils.SendOutput(err, w, r, data, "UpdateFaq")
}
func DeleteFaq(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	err := service.DeleteFaq(id)
	utils.SendOutput(err, w, r, "Deleted successfully", "DeleteFaq")
}

func GetAllFaq(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	data, err := service.GetAllFaq()
	utils.SendOutput(err, w, r, data, "GetAllFaq")
}
func parseFaq(r *http.Request) (entity.FAQDB, error) {
	var faq entity.FAQDB
	err := json.NewDecoder(r.Body).Decode(&faq)
	if err != nil {
		return entity.FAQDB{}, err
	}
	return faq, nil
}
