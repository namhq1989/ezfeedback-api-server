package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type ProjectSettingRepository interface {
	Create(ctx *appcontext.AppContext, setting ProjectSetting) error
	Update(ctx *appcontext.AppContext, setting ProjectSetting) error
	FindByProjectID(ctx *appcontext.AppContext, projectID string) (*ProjectSetting, error)
}

type ProjectSetting struct {
	ID                     string
	ProjectID              string
	IsFeedbackPublic       bool
	AllowAnonymousFeedback bool
	EnableVoting           bool
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func NewProjectSetting(projectID string, isFeedbackPublic bool, allowAnonymousFeedback bool, enableVoting bool) (*ProjectSetting, error) {
	var (
		now = manipulation.NowUTC()
	)

	var s = &ProjectSetting{
		ID:                     uuid.New(),
		IsFeedbackPublic:       isFeedbackPublic,
		AllowAnonymousFeedback: allowAnonymousFeedback,
		EnableVoting:           enableVoting,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if err := s.SetProjectID(projectID); err != nil {
		return nil, err
	}

	return s, nil
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
