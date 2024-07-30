package notify

import (
	"futureEVChronJobs/pkg/entity"
	notificationrepo "futureEVChronJobs/pkg/repo/notification"
	utils "futureEVChronJobs/pkg/util"
)

var repo = notificationrepo.NewRepository("notifications")

type str struct{}

func NewService() Notify {
	return &str{}
}

// SendNotification implements Notify.
func (s *str) SendNotification(title string, body string, userId string, ntype, token string) error {
	utils.SendNotification(title, body, token)
	_, err := repo.InsertOne(entity.Notification{
		Title:  title,
		Body:   body,
		UserId: userId,
		Token:  token,
		Type:   ntype,
	})
	return err
}
