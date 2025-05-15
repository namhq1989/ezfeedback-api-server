package dto

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
)

type FeedbackStateHistory struct {
	ID        string                    `json:"id"`
	User      User                      `json:"user"`
	State     string                    `json:"state"`
	CreatedAt *httprespond.TimeResponse `json:"createdAt"`
}

func (FeedbackStateHistory) FromDomain(history domain.FeedbackStateHistory, user domain.User) FeedbackStateHistory {
	return FeedbackStateHistory{
		ID:        history.ID,
		User:      User{}.FromDomain(user),
		State:     history.State.String(),
		CreatedAt: httprespond.NewTimeResponse(history.CreatedAt),
	}
}
