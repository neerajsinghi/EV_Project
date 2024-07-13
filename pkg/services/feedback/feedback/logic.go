package feedback

import (
	"bikeRental/pkg/entity"
	feedbackrepo "bikeRental/pkg/repo/feedback"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = feedbackrepo.NewRepository("feedback")

type feed struct{}

func New() FeedbackI {
	return &feed{}
}

func (f *feed) AddFeedback(feedback entity.Feedback) (string, error) {
	return repo.InsertOne(feedback)
}

func (f *feed) GetFeedbacks() ([]entity.Feedback, error) {
	return repo.Find(nil, nil)
}

func (f *feed) DeleteFeedback(feedbackID string) error {
	idObj, _ := primitive.ObjectIDFromHex(feedbackID)
	return repo.DeleteOne(bson.M{"_id": idObj})
}
