package feedback

import (
	"encoding/json"
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/services/feedback/feedback"
	utils "futureEVChronJobs/pkg/util"
	"net/http"

	"github.com/gorilla/mux"
)

var feed = feedback.New()

// AddFeedback adds a feedback to the database
func AddFeedback(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	feedback, err := parseFeed(r)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(entity.Response{Status: false, Error: "Invalid request"})
		return
	}
	id, err := feed.AddFeedback(feedback)
	utils.SendOutput(err, w, r, id, feedback, "AddFeedback")
}

// GetFeedbacks returns all feedbacks
func GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	data, err := feed.GetFeedbacks()
	utils.SendOutput(err, w, r, data, nil, "GetFeedbacks")
}

// DeleteFeedback deletes a feedback
func DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	err := feed.DeleteFeedback(id)
	utils.SendOutput(err, w, r, "Deleted successfully", nil, "DeleteFeedback")
}

func parseFeed(r *http.Request) (entity.Feedback, error) {
	var feedback entity.Feedback
	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		return entity.Feedback{}, err
	}
	return feedback, nil
}
