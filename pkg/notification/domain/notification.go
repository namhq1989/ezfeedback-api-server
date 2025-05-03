package domain

import (
	"time"

	"github.com/namhq1989/go-utilities/appcontext"
)

type NotificationRepository interface {
	Create(ctx *appcontext.AppContext, notification Notification) error
	FindWithFilter(ctx *appcontext.AppContext, filter NotificationFilter) ([]Notification, error)
	CleanupStale(ctx *appcontext.AppContext) error
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
