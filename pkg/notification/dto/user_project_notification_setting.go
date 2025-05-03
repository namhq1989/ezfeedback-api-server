package dto

import "github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"

type UserProjectNotificationSetting struct {
	ID                   string `json:"id"`
	ReceiveNewFeedback   bool   `json:"receiveNewFeedback"`
	ReceiveDailySummary  bool   `json:"receiveDailySummary"`
	ReceiveWeeklySummary bool   `json:"receiveWeeklySummary"`
}

func (UserProjectNotificationSetting) FromDomain(setting domain.UserProjectNotificationSetting) UserProjectNotificationSetting {
	return UserProjectNotificationSetting{
		ID:                   setting.ID,
		ReceiveNewFeedback:   setting.ReceiveNewFeedback,
		ReceiveDailySummary:  setting.ReceiveDailySummary,
		ReceiveWeeklySummary: setting.ReceiveWeeklySummary,
	}
}
