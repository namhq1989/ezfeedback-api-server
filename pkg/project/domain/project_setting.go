package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/uuid"
)

type ProjectSetting struct {
	ID                     string
	ProjectID              string
	IsFeedbackPublic       bool
	AllowAnonymousFeedback bool
	EnableVoting           bool
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func NewProjectSetting(projectID string) *ProjectSetting {
	var s = &ProjectSetting{
		IsFeedbackPublic:       false,
		AllowAnonymousFeedback: true,
		EnableVoting:           true,
	}

	if err := s.SetProjectID(projectID); err != nil {
		return nil
	}

	return s
}

func (s *ProjectSetting) SetProjectID(projectID string) error {
	if !uuid.IsValidID(projectID) {
		return apperrors.Project.InvalidProjectID
	}

	s.ProjectID = projectID
	s.SetUpdatedAt()
	return nil
}

func (s *ProjectSetting) SetIsFeedbackPublic(isFeedbackPublic bool) {
	s.IsFeedbackPublic = isFeedbackPublic
	s.SetUpdatedAt()
}

func (s *ProjectSetting) SetAllowAnonymousFeedback(allowAnonymousFeedback bool) {
	s.AllowAnonymousFeedback = allowAnonymousFeedback
	s.SetUpdatedAt()
}

func (s *ProjectSetting) SetEnableVoting(enableVoting bool) {
	s.EnableVoting = enableVoting
	s.SetUpdatedAt()
}

func (s *ProjectSetting) SetUpdatedAt() {
	s.UpdatedAt = manipulation.NowUTC()
}
