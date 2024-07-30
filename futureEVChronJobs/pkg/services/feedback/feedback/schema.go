package feedback

import "futureEVChronJobs/pkg/entity"

type FeedbackI interface {
	AddFeedback(feedback entity.Feedback) (string, error)
	GetFeedbacks() ([]entity.FeedbackOut, error)
	DeleteFeedback(feedbackID string) error
}
