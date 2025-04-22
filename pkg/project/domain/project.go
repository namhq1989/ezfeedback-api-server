package domain

import (
	"fmt"
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type ProjectRepository interface {
	Create(ctx *appcontext.AppContext, project Project) error
	Update(ctx *appcontext.AppContext, project Project) error
	FindByID(ctx *appcontext.AppContext, projectID string) (*Project, error)
	FindByUserID(ctx *appcontext.AppContext, userID string) ([]Project, error)
}

const (
	projectSlugSuffixLength = 6
)

type Project struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Slug        string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProject(userID, title, description string) (*Project, error) {
	var (
		now = manipulation.NowUTC()
	)

	var p = &Project{
		ID:        uuid.New(),
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := p.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := p.SetTitle(title); err != nil {
		return nil, err
	}
	if err := p.SetDescription(description); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Project) SetUserID(userID string) error {
	if !uuid.IsValidID(userID) {
		return apperrors.User.InvalidUserID
	}

	p.UserID = userID
	p.SetUpdatedAt()
	return nil
}

func (p *Project) SetTitle(title string) error {
	if title == "" || len(title) < 3 || len(title) > 255 {
		return apperrors.Common.InvalidTitle
	}

	p.Title = title
	p.Slug = fmt.Sprintf("%s-%s", manipulation.Slugify(title), manipulation.RandomAlphaNumeric(projectSlugSuffixLength))
	p.SetUpdatedAt()
	return nil
}

func (p *Project) SetDescription(description string) error {
	if len(description) > 2000 {
		return apperrors.Common.InvalidDescription
	}

	p.Description = description
	p.SetUpdatedAt()
	return nil
}

func (p *Project) SetUpdatedAt() {
	p.UpdatedAt = manipulation.NowUTC()
}
