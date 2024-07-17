package notify

import (
	"bikeRental/pkg/entity"
	notificationrepo "bikeRental/pkg/repo/notification"
	utils "bikeRental/pkg/util"
	"errors"
)

var repo = notificationrepo.NewRepository("notifications")

type str struct{}

func NewService() Notify {
	return &str{}
}

// SendNotification implements Notify.
func (s *str) SendNotification(title string, body string, userId string, token, ntype string) error {
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
func (s *str) SendMultipleNotifications(title, body, ntype string, userIds []string, tokens []string) error {
	if len(userIds) != len(tokens) {
		return errors.New("userIds and tokens must have the same length")
	}
	for i := range userIds {
		utils.SendNotification(title, body, tokens[i])
		_, err := repo.InsertOne(entity.Notification{
			Title:  title,
			Body:   body,
			UserId: userIds[i],
			Token:  tokens[i],
			Type:   ntype,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllNotifications implements Notify.
func (s *str) GetAllNotifications() ([]entity.Notification, error) {
	return repo.Find(nil, nil)
}
