package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/invitation/domain"
)

type UserInvitationMapper struct{}

func (UserInvitationMapper) FromModelToDomain(invitation model.UserInvitations) (*domain.UserInvitation, error) {
	var result = &domain.UserInvitation{
		ID:        invitation.ID,
		Email:     invitation.Email,
		InviterID: invitation.InviterID,
		ProjectID: invitation.ProjectID,
		Role:      domain.ToProjectRole(invitation.Role.String()),
		Status:    domain.ToUserInvitationStatus(invitation.Status.String()),
		ExpiresAt: invitation.ExpiresAt,
		CreatedAt: invitation.CreatedAt,
		UpdatedAt: invitation.UpdatedAt,
	}

	return result, nil
}

func (UserInvitationMapper) FromDomainToModel(invitation domain.UserInvitation) (*model.UserInvitations, error) {
	var result = &model.UserInvitations{
		ID:        invitation.ID,
		Email:     invitation.Email,
		InviterID: invitation.InviterID,
		ProjectID: invitation.ProjectID,
		Role:      model.ProjectRole(invitation.Role.String()),
		Status:    model.InvitationStatus(invitation.Status.String()),
		ExpiresAt: invitation.ExpiresAt,
		CreatedAt: invitation.CreatedAt,
		UpdatedAt: invitation.UpdatedAt,
	}

	return result, nil
}
