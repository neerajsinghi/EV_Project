package notify

type Notify interface {
	SendNotification(title, body, userId, ntype, token string) error
}
