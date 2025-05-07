package shared

import "github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"

type Service struct {
	feedbackRepository    domain.FeedbackRepository
	cachingRepository     domain.CachingRepository
	externalAPIRepository domain.ExternalAPIRepository
}

func NewService(
	feedbackRepository domain.FeedbackRepository,
	cachingRepository domain.CachingRepository,
	externalAPIRepository domain.ExternalAPIRepository,
) Service {
	return Service{
		feedbackRepository:    feedbackRepository,
		cachingRepository:     cachingRepository,
		externalAPIRepository: externalAPIRepository,
	}
}
