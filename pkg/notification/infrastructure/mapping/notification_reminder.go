package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
)

type NotificationReminderMapper struct{}

func (NotificationReminderMapper) FromModelToDomain(reminder model.NotificationReminders) (*domain.NotificationReminder, error) {
	var result = &domain.NotificationReminder{
		ID:        reminder.ID,
		ProjectID: reminder.ProjectID,
		CreatedAt: reminder.CreatedAt,
	}

	return result, nil
}

func (NotificationReminderMapper) FromDomainToModel(reminder domain.NotificationReminder) (*model.NotificationReminders, error) {
	var result = &model.NotificationReminders{
		ID:        reminder.ID,
		ProjectID: reminder.ProjectID,
		CreatedAt: reminder.CreatedAt,
	}

	return result, nil
}
