package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type ProjectInvitationRepository interface {
	Create(ctx *appcontext.AppContext, invitation ProjectInvitation) error
	Update(ctx *appcontext.AppContext, invitation ProjectInvitation) error
	Delete(ctx *appcontext.AppContext, invitationID string) error
	FindByID(ctx *appcontext.AppContext, invitationID string) (*ProjectInvitation, error)
	FindByProjectID(ctx *appcontext.AppContext, projectID string) ([]ProjectInvitation, error)
	FindByEmail(ctx *appcontext.AppContext, email string) ([]ProjectInvitation, error)
	MarkExpiredInvitations(ctx *appcontext.AppContext) error
	CleanupStale(ctx *appcontext.AppContext) error
	CountPendingByProjectID(ctx *appcontext.AppContext, projectID string) (int64, error)
}

const (
	invitationTTL = 7 * 24 * time.Hour
)

type ProjectInvitation struct {
	ID        string
	Email     string
	InviterID string
	ProjectID string
	Role      ProjectRole
	Status    ProjectInvitationStatus
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProjectInvitation(email, inviterID, projectID, role string) (*ProjectInvitation, error) {
	var (
		now = manipulation.NowUTC()
	)

	var i = &ProjectInvitation{
		ID:        uuid.New(),
		Status:    ProjectInvitationStatusPending,
		ExpiresAt: now.Add(invitationTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := i.SetEmail(email); err != nil {
		return nil, err
	}
	if err := i.SetInviterID(inviterID); err != nil {
		return nil, err
	}
	if err := i.SetProjectID(projectID); err != nil {
		return nil, err
	}
	if err := i.SetRole(role); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *ProjectInvitation) SetEmail(email string) error {
	if !validation.IsValidEmail(email) {
		return apperrors.Common.InvalidEmail
	}

	i.Email = email
	i.SetUpdatedAt()
	return nil
}

func (i *ProjectInvitation) SetInviterID(inviterID string) error {
	if !uuid.IsValidID(inviterID) {
		return apperrors.Invitation.InvalidInviterID
	}

	i.InviterID = inviterID
	i.SetUpdatedAt()
	return nil
}

func (i *ProjectInvitation) SetProjectID(projectID string) error {
	if !uuid.IsValidID(projectID) {
		return apperrors.Project.InvalidProjectID
	}

	i.ProjectID = projectID
	i.SetUpdatedAt()
	return nil
}

func (i *ProjectInvitation) SetRole(role string) error {
	dRole := ToProjectRole(role)
	if !dRole.IsValid() {
		return apperrors.Project.InvalidRole
	}

	i.Role = dRole
	i.SetUpdatedAt()
	return nil
}

func (i *ProjectInvitation) SetStatus(status string) error {
	dStatus := ToProjectInvitationStatus(status)
	if !dStatus.IsValid() {
		return apperrors.Common.InvalidStatus
	}

	i.Status = dStatus
	i.SetUpdatedAt()
	return nil
}

func (i *ProjectInvitation) SetUpdatedAt() {
	i.UpdatedAt = manipulation.NowUTC()
}

func (i *ProjectInvitation) IsExpired() bool {
	return i.ExpiresAt.Before(manipulation.NowUTC())
}
