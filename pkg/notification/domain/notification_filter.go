package domain

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/go-utilities/uuid"
)

const (
	notificationQueryLimit int64 = 20
)

type NotificationFilter struct {
	UserID string
	Page   int64
	Limit  int64
}

func NewNotificationFilter(userID string, page int64) (*NotificationFilter, error) {
	if !uuid.IsValidID(userID) {
		return nil, apperrors.User.InvalidUserID
	}

	return &NotificationFilter{
		UserID: userID,
		Page:   page,
		Limit:  notificationQueryLimit,
	}, nil
}

func IsEndOfNotificationResult(total int64) bool {
	return total < notificationQueryLimit
}
