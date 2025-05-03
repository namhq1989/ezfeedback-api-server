package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
)

type UserProjectNotificationSettingMapper struct{}

func (UserProjectNotificationSettingMapper) FromModelToDomain(setting model.UserProjectNotificationSettings) (*domain.UserProjectNotificationSetting, error) {
	var result = &domain.UserProjectNotificationSetting{
		ID:                   setting.ID,
		UserID:               setting.UserID,
		ProjectID:            setting.ProjectID,
		ReceiveNewFeedback:   setting.ReceiveNewFeedback,
		ReceiveDailySummary:  setting.ReceiveDailySummary,
		ReceiveWeeklySummary: setting.ReceiveWeeklySummary,
		CreatedAt:            setting.CreatedAt,
		UpdatedAt:            setting.UpdatedAt,
	}

	return result, nil
}

func (UserProjectNotificationSettingMapper) FromDomainToModel(setting domain.UserProjectNotificationSetting) (*model.UserProjectNotificationSettings, error) {
	var result = &model.UserProjectNotificationSettings{
		ID:                   setting.ID,
		UserID:               setting.UserID,
		ProjectID:            setting.ProjectID,
		ReceiveNewFeedback:   setting.ReceiveNewFeedback,
		ReceiveDailySummary:  setting.ReceiveDailySummary,
		ReceiveWeeklySummary: setting.ReceiveWeeklySummary,
		CreatedAt:            setting.CreatedAt,
		UpdatedAt:            setting.UpdatedAt,
	}

	return result, nil
}
