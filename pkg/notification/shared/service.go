package shared

import "github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"

type Service struct {
	userProjectNotificationSettingRepository domain.UserProjectNotificationSettingRepository
	cachingRepository                        domain.CachingRepository
}

func NewService(
	userProjectNotificationSettingRepository domain.UserProjectNotificationSettingRepository,
	cachingRepository domain.CachingRepository,
) Service {
	return Service{
		userProjectNotificationSettingRepository: userProjectNotificationSettingRepository,
		cachingRepository:                        cachingRepository,
	}
}
