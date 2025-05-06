package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type NotificationRepository interface {
	Create(ctx *appcontext.AppContext, notification Notification) error
	FindWithFilter(ctx *appcontext.AppContext, filter NotificationFilter) ([]Notification, error)
	CountWithFilter(ctx *appcontext.AppContext, filter NotificationFilter) (int64, error)
	CleanupStale(ctx *appcontext.AppContext) error

	GenerateNewFeedbackContent(_ *appcontext.AppContext, language string, projectTitle string) string
}

type Notification struct {
	ID        string
	UserID    string
	Type      NotificationType
	IsRead    bool
	Metadata  NotificationMetadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFeedbackNotification(userID string, metadata NotificationMetadata) (*Notification, error) {
	var (
		now = manipulation.NowUTC()
	)

	var n = &Notification{
		ID:        uuid.New(),
		Type:      NotificationTypeNewFeedback,
		IsRead:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := n.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := n.SetMetadata(metadata); err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Notification) SetUserID(userID string) error {
	if !uuid.IsValidID(userID) {
		return apperrors.User.InvalidUserID
	}

	n.UserID = userID
	return nil
}

func (n *Notification) SetMetadata(metadata NotificationMetadata) error {
	if n.Type.IsNewFeedback() {
		if !uuid.IsValidID(metadata.ProjectID) || metadata.ProjectTitle == "" {
			return apperrors.Notification.InvalidMetadata
		}
	}

	n.Metadata = metadata
	return nil
}
