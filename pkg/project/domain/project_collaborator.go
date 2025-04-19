package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/uuid"
)

type ProjectCollaborator struct {
	ID        string
	ProjectID string
	UserID    string
	Role      ProjectRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProjectCollaborator(projectID, userID, role string) (*ProjectCollaborator, error) {
	var (
		now = manipulation.NowUTC()
	)

	var c = &ProjectCollaborator{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.SetProjectID(projectID); err != nil {
		return nil, err
	}
	if err := c.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := c.SetRole(role); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *ProjectCollaborator) SetProjectID(projectID string) error {
	if !uuid.IsValidID(projectID) {
		return apperrors.Project.InvalidProjectID
	}
	c.ProjectID = projectID
	return nil
}

func (c *ProjectCollaborator) SetUserID(userID string) error {
	if !uuid.IsValidID(userID) {
		return apperrors.User.InvalidUserID
	}
	c.UserID = userID
	return nil
}

func (c *ProjectCollaborator) SetRole(role string) error {
	var dRole = ToProjectRole(role)
	if !dRole.IsValid() {
		return apperrors.Project.InvalidRole
	}

	c.Role = ProjectRole(role)
	return nil
}
