package domain

type NotificationType string

const (
	NotificationTypeUnknown     NotificationType = ""
	NotificationTypeNewFeedback NotificationType = "new_feedback"
	NotificationTypeSystem      NotificationType = "system"
)

func (n NotificationType) String() string {
	return string(n)
}

func (n NotificationType) IsValid() bool {
	return n != NotificationTypeUnknown
}

func (n NotificationType) IsNewFeedback() bool {
	return n == NotificationTypeNewFeedback
}

func (n NotificationType) IsSystem() bool {
	return n == NotificationTypeSystem
}

func ToNotificationType(s string) NotificationType {
	switch s {
	case NotificationTypeNewFeedback.String():
		return NotificationTypeNewFeedback
	case NotificationTypeSystem.String():
		return NotificationTypeSystem
	default:
		return NotificationTypeUnknown
	}
}
