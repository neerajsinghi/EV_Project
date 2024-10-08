package feedback

import (
	"bikeRental/pkg/entity"
	feedbackrepo "bikeRental/pkg/repo/feedback"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = feedbackrepo.NewRepository("feedback")

type feed struct{}

func New() FeedbackI {
	return &feed{}
}

func (f *feed) AddFeedback(feedback entity.Feedback) (string, error) {
	data, err := repo.InsertOne(feedback)
	if err != nil {
		return "", errors.New("error in inserting feedback")
	}
	return data, nil
}

func (f *feed) GetFeedbacks() ([]entity.FeedbackOut, error) {
	pipeline := bson.A{
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "userObjectId", Value: bson.D{{Key: "$toObjectId", Value: "$profile_id"}}},
					{Key: "bookingObjectId", Value: bson.D{{Key: "$toObjectId", Value: "$booking_id"}}},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "booking"},
					{Key: "localField", Value: "bookingObjectId"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "booking"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$booking"}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "users"},
					{Key: "localField", Value: "userObjectId"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "profile"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "feedback", Value: 1},
					{Key: "ratings", Value: 1},
					{Key: "profile", Value: 1},
					{Key: "booking", Value: 1},
				},
			},
		},
	}
	data, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, errors.New("error in finding feedback")
	}
	return data, nil

}

func (f *feed) DeleteFeedback(feedbackID string) error {
	idObj, _ := primitive.ObjectIDFromHex(feedbackID)
	err := repo.DeleteOne(bson.M{"_id": idObj})
	if err != nil {
		return errors.New("error in deleting feedback")
	}
	return nil
}
