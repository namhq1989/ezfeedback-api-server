package mapping

import (
	"github.com/goccy/go-json"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
)

type NotificationMapper struct{}

func (NotificationMapper) FromModelToDomain(notification model.Notifications) (*domain.Notification, error) {
	var metadata NotificationMetadata
	if notification.Metadata != "" {
		if err := json.Unmarshal([]byte(notification.Metadata), &metadata); err != nil {
			return nil, err
		}
	}

	var result = &domain.Notification{
		ID:     notification.ID,
		UserID: notification.UserID,
		Type:   domain.ToNotificationType(notification.Type.String()),
		IsRead: notification.IsRead,
		Metadata: domain.NotificationMetadata{
			ProjectID:    metadata.ProjectID,
			ProjectTitle: metadata.ProjectTitle,
		},
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}

	return result, nil
}

func (NotificationMapper) FromDomainToModel(notification domain.Notification) (*model.Notifications, error) {
	var result = &model.Notifications{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Type:      model.NotificationType(notification.Type.String()),
		IsRead:    notification.IsRead,
		Metadata:  "",
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}

	var metadata = NotificationMetadata{
		ProjectID:    notification.Metadata.ProjectID,
		ProjectTitle: notification.Metadata.ProjectTitle,
	}
	if metadataBytes, err := json.Marshal(metadata); err != nil {
		return nil, err
	} else {
		result.Metadata = string(metadataBytes)
	}

	return result, nil
}
