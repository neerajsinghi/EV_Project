package feedback

import "bikeRental/pkg/entity"

type FeedbackI interface {
	AddFeedback(feedback entity.Feedback) (string, error)
	GetFeedbacks() ([]entity.Feedback, error)
	DeleteFeedback(feedbackID string) error
}
