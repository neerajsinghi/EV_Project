package notify

import "bikeRental/pkg/entity"

type Notify interface {
	SendNotification(title, body, userId, ntype, token string) error
	SendMultipleNotifications(title, body, ntype string, userIds []string, token []string) error
	GetAllNotifications() ([]entity.Notification, error)
}
