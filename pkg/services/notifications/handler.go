package notifications

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/notifications/notify"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

var service = notify.NewService()

// Notify defines the required methods for the notification service.
func SendNotification(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)

	notification, err := parseNotificationRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "unable to decode request body"})
		return
	}

	err = service.SendNotification(notification.Title, notification.Body, notification.UserId, notification.Type, notification.Token)
	utils.SendOutput(err, w, r, "Notification sent successfully", "SendNotification")
}

// GetAllNotifications returns all notifications.
func GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)

	notifications, err := service.GetAllNotifications()
	utils.SendOutput(err, w, r, notifications, "GetAllNotifications")
}
func parseNotificationRequest(r *http.Request) (entity.Notification, error) {
	var notification entity.Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		return notification, errors.Wrapf(err, "unable to decode request body")
	}
	return notification, nil
}
func parseMultiNotificationRequest(r *http.Request) (entity.NotificationMulti, error) {
	var notification entity.NotificationMulti
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		return notification, errors.Wrapf(err, "unable to decode request body")
	}
	return notification, nil
}

// SendMultipleNotifications sends notifications to multiple users.
func SendMultipleNotifications(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)

	notification, err := parseMultiNotificationRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "unable to decode request body"})
		return
	}

	err = service.SendMultipleNotifications(notification.Title, notification.Body, notification.Type, notification.UserIds, notification.Tokens)
	utils.SendOutput(err, w, r, "Notifications sent successfully", "SendMultipleNotifications")
}
