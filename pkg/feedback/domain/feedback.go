package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/uuid"
)

type Feedback struct {
	ID          string
	ProjectID   string
	UserID      *string
	CategoryID  string
	Content     string
	Rating      int32
	IsAnonymous bool
	State       FeedbackState
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewFeedback(projectID string, userID *string, categoryID, content string, rating int32) (*Feedback, error) {
	var (
		now = manipulation.NowUTC()
	)
	var f = &Feedback{
		ID:        uuid.New(),
		State:     FeedbackStateNew,
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := f.SetProjectID(projectID); err != nil {
		return nil, err
	}
	if err := f.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := f.SetCategoryID(categoryID); err != nil {
		return nil, err
	}
	if err := f.SetContent(content); err != nil {
		return nil, err
	}
	if err := f.SetRating(rating); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *Feedback) SetProjectID(projectID string) error {
	if !uuid.IsValidID(projectID) {
		return apperrors.Project.InvalidProjectID
	}

	f.ProjectID = projectID
	f.SetUpdatedAt()
	return nil
}

func (f *Feedback) SetUserID(userID *string) error {
	if userID != nil && !uuid.IsValidID(*userID) {
		return apperrors.User.InvalidUserID
	}

	f.UserID = userID
	f.IsAnonymous = userID == nil
	f.SetUpdatedAt()
	return nil
}

func (f *Feedback) SetCategoryID(categoryID string) error {
	if !uuid.IsValidID(categoryID) {
		return apperrors.Project.InvalidCategory
	}

	f.CategoryID = categoryID
	f.SetUpdatedAt()
	return nil
}

func (f *Feedback) SetContent(content string) error {
	if len(content) > 2000 {
		return apperrors.Feedback.InvalidContent
	}

	f.Content = content
	f.SetUpdatedAt()
	return nil
}

func (f *Feedback) SetRating(rating int32) error {
	if rating < 1 || rating > 5 {
		return apperrors.Feedback.InvalidRating
	}

	f.Rating = rating
	f.SetUpdatedAt()
	return nil
}

func (f *Feedback) SetState(state string) error {
	var dState = ToFeedbackState(state)
	if !dState.IsValid() {
		return apperrors.Feedback.InvalidState
	}

	f.State = dState
	f.SetUpdatedAt()
	return nil
}

func (f *Feedback) SetStatus(status string) error {
	var dStatus = ToStatus(status)
	if !dStatus.IsValid() {
		return apperrors.Common.InvalidStatus
	}

	f.Status = dStatus
	f.SetUpdatedAt()
	return nil
}

func (f *Feedback) SetUpdatedAt() {
	f.UpdatedAt = manipulation.NowUTC()
}
