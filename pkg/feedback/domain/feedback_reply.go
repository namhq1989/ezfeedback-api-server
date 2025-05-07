package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type FeedbackReplyRepository interface {
	Create(ctx *appcontext.AppContext, reply FeedbackReply) error
	Update(ctx *appcontext.AppContext, reply FeedbackReply) error
	Delete(ctx *appcontext.AppContext, reply FeedbackReply) error
	FindByID(ctx *appcontext.AppContext, replyID string) (*FeedbackReply, error)
	FindWithFilter(ctx *appcontext.AppContext, filter FeedbackReplyFilter) ([]FeedbackReply, error)
}

type FeedbackReply struct {
	ID         string
	FeedbackID string
	UserID     string
	Content    string
	IsEdited   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewFeedbackReply(feedbackID, userID, content string) (*FeedbackReply, error) {
	var now = manipulation.NowUTC()
	var r = &FeedbackReply{
		ID:        uuid.New(),
		IsEdited:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := r.SetFeedbackID(feedbackID); err != nil {
		return nil, err
	}
	if err := r.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := r.SetContent(content); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *FeedbackReply) SetFeedbackID(feedbackID string) error {
	if !uuid.IsValidID(feedbackID) {
		return apperrors.Feedback.InvalidFeedbackID
	}

	r.FeedbackID = feedbackID
	r.SetUpdatedAt()
	return nil
}

func (r *FeedbackReply) SetUserID(userID string) error {
	if !uuid.IsValidID(userID) {
		return apperrors.User.InvalidUserID
	}

	r.UserID = userID
	r.SetUpdatedAt()
	return nil
}

func (r *FeedbackReply) SetContent(content string) error {
	if content == "" || len(content) > 2000 {
		return apperrors.Feedback.InvalidContent
	}

	r.Content = content
	r.SetUpdatedAt()
	return nil
}

func (r *FeedbackReply) SetUpdatedAt() {
	r.UpdatedAt = manipulation.NowUTC()
}
