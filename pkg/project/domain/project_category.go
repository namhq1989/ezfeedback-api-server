package domain

import (
	"fmt"
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/uuid"
)

const (
	projectCategorySlugSuffixLength = 4
)

type ProjectCategory struct {
	ID        string
	ProjectID string
	Name      string
	Slug      string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProjectCategory(projectID, name string) (*ProjectCategory, error) {
	var (
		now = manipulation.NowUTC()
	)

	var c = &ProjectCategory{
		ID:        uuid.New(),
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.SetProjectID(projectID); err != nil {
		return nil, err
	}
	if err := c.SetName(name); err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *ProjectCategory) SetProjectID(projectID string) error {
	if !uuid.IsValidID(projectID) {
		return apperrors.Project.InvalidProjectID
	}
	c.ProjectID = projectID
	c.SetUpdatedAt()
	return nil
}

func (c *ProjectCategory) SetName(name string) error {
	if name == "" || len(name) < 3 || len(name) > 255 {
		return apperrors.Common.InvalidName
	}

	c.Name = name
	c.Slug = fmt.Sprintf("%s-%s", manipulation.Slugify(name), manipulation.RandomAlphaNumeric(projectCategorySlugSuffixLength))
	c.SetUpdatedAt()
	return nil
}

func (c *ProjectCategory) SetStatus(status string) error {
	var dStatus = ToStatus(status)
	if !dStatus.IsValid() {
		return apperrors.Common.InvalidStatus
	}

	c.Status = dStatus
	c.SetUpdatedAt()
	return nil
}

func (c *ProjectCategory) SetUpdatedAt() {
	c.UpdatedAt = manipulation.NowUTC()
}
