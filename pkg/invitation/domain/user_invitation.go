package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/go-utilities/uuid"
)

const (
	invitationTTL = 24 * time.Hour
)

type UserInvitation struct {
	ID        string
	Email     string
	InviterID string
	ProjectID string
	Role      ProjectRole
	Token     string
	Status    UserInvitationStatus
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserInvitation(email, inviterID, projectID, role, token string) (*UserInvitation, error) {
	var (
		now = manipulation.NowUTC()
	)

	var i = &UserInvitation{
		ID:        uuid.New(),
		Status:    UserInvitationStatusPending,
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
	if err := i.SetToken(token); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *UserInvitation) SetEmail(email string) error {
	if !validation.IsValidEmail(email) {
		return apperrors.Common.InvalidEmail
	}

	i.Email = email
	i.SetUpdatedAt()
	return nil
}

func (i *UserInvitation) SetInviterID(inviterID string) error {
	if !uuid.IsValidID(inviterID) {
		return apperrors.Invitation.InvalidInviterID
	}

	i.InviterID = inviterID
	i.SetUpdatedAt()
	return nil
}

func (i *UserInvitation) SetProjectID(projectID string) error {
	if !uuid.IsValidID(projectID) {
		return apperrors.Project.InvalidProjectID
	}

	i.ProjectID = projectID
	i.SetUpdatedAt()
	return nil
}

func (i *UserInvitation) SetRole(role string) error {
	dRole := ToProjectRole(role)
	if !dRole.IsValid() {
		return apperrors.Project.InvalidRole
	}

	i.Role = dRole
	i.SetUpdatedAt()
	return nil
}

func (i *UserInvitation) SetToken(token string) error {
	i.Token = token
	i.SetUpdatedAt()
	return nil
}

func (i *UserInvitation) SetStatus(status string) error {
	dStatus := ToUserInvitationStatus(status)
	if !dStatus.IsValid() {
		return apperrors.Common.InvalidStatus
	}

	i.Status = dStatus
	i.SetUpdatedAt()
	return nil
}

func (i *UserInvitation) SetUpdatedAt() {
	i.UpdatedAt = manipulation.NowUTC()
}

func (i *UserInvitation) IsValidToken(token string) bool {
	return i.Token == token
}

func (i *UserInvitation) IsExpired() bool {
	return i.ExpiresAt.Before(manipulation.NowUTC())
}
