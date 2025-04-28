package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

const (
	projectSettingDefaultColor = "#2563eb" // shadcn blue-600
)

type ProjectSettingRepository interface {
	Create(ctx *appcontext.AppContext, setting ProjectSetting) error
	Update(ctx *appcontext.AppContext, setting ProjectSetting) error
	FindByProjectID(ctx *appcontext.AppContext, projectID string) (*ProjectSetting, error)
}

type ProjectSetting struct {
	ID           string
	ProjectID    string
	Domain       string
	PrimaryColor string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewProjectSetting(projectID, domain, primaryColor string) (*ProjectSetting, error) {
	var (
		now = manipulation.NowUTC()
	)

	var s = &ProjectSetting{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.SetProjectID(projectID); err != nil {
		return nil, err
	}
	if err := s.SetDomain(domain); err != nil {
		return nil, err
	}
	if err := s.SetPrimaryColor(primaryColor); err != nil {
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

func (s *ProjectSetting) SetDomain(domain string) error {
	if len(domain) > 0 && !validation.IsValidDomain(domain) {
		return apperrors.Project.InvalidDomain
	}

	s.Domain = domain
	s.SetUpdatedAt()
	return nil
}

func (s *ProjectSetting) SetPrimaryColor(primaryColor string) error {
	if len(primaryColor) > 0 && !validation.IsValidHexColor(primaryColor) {
		return apperrors.Project.InvalidPrimaryColor
	}

	s.PrimaryColor = primaryColor
	if s.PrimaryColor == "" {
		s.PrimaryColor = projectSettingDefaultColor
	}
	s.SetUpdatedAt()
	return nil
}

func (s *ProjectSetting) SetUpdatedAt() {
	s.UpdatedAt = manipulation.NowUTC()
}
