package dto

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
)

type FeedbackReply struct {
	ID        string                    `json:"id"`
	User      User                      `json:"user"`
	Content   string                    `json:"content"`
	IsEdited  bool                      `json:"isEdited"`
	CreatedAt *httprespond.TimeResponse `json:"createdAt"`
	UpdatedAt *httprespond.TimeResponse `json:"updatedAt"`
}

func (FeedbackReply) FromDomain(reply domain.FeedbackReply, user domain.User) FeedbackReply {
	return FeedbackReply{
		ID:        reply.ID,
		User:      User{}.FromDomain(user),
		Content:   reply.Content,
		IsEdited:  reply.IsEdited,
		CreatedAt: httprespond.NewTimeResponse(reply.CreatedAt),
		UpdatedAt: httprespond.NewTimeResponse(reply.UpdatedAt),
	}
}
