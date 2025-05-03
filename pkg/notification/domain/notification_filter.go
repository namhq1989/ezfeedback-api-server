package domain

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/pagetoken"
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

func NewNotificationFilter(pageToken, userID string) (*NotificationFilter, error) {
	if !uuid.IsValidID(userID) {
		return nil, apperrors.User.InvalidUserID
	}

	pt := pagetoken.Decode(pageToken)

	return &NotificationFilter{
		UserID: userID,
		Page:   pt.Page,
		Limit:  notificationQueryLimit,
	}, nil
}

func IsEndOfNotificationResult(total int64) bool {
	return total < notificationQueryLimit
}
