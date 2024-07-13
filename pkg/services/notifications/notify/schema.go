package notify

import "bikeRental/pkg/entity"

type Notify interface {
	SendNotification(title string, body string, userId string, token string) error
	SendMultipleNotifications(title string, body string, userIds []string, token []string) error
	GetAllNotifications() ([]entity.Notification, error)
}
