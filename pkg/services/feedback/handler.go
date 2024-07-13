package feedback

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/feedback/feedback"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var feed = feedback.New()

// AddFeedback adds a feedback to the database
func AddFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	feedback, err := parseFeed(r)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(entity.Response{Status: false, Error: "Invalid request"})
		return
	}
	id, err := feed.AddFeedback(feedback)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entity.Response{Status: false, Error: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.Response{Status: true, Data: id})
}

// GetFeedbacks returns all feedbacks
func GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, err := feed.GetFeedbacks()
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entity.Response{Status: false, Error: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.Response{Status: true, Data: data})
}

// DeleteFeedback deletes a feedback
func DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := mux.Vars(r)["id"]
	err := feed.DeleteFeedback(id)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entity.Response{Status: false, Error: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.Response{Status: true})
}

func parseFeed(r *http.Request) (entity.Feedback, error) {
	var feedback entity.Feedback
	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		return entity.Feedback{}, err
	}
	return feedback, nil
}
