package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectInvitationMapper struct{}

func (ProjectInvitationMapper) FromModelToDomain(invitation model.ProjectInvitations) (*domain.ProjectInvitation, error) {
	var result = &domain.ProjectInvitation{
		ID:        invitation.ID,
		Email:     invitation.Email,
		InviterID: invitation.InviterID,
		ProjectID: invitation.ProjectID,
		Role:      domain.ToProjectRole(invitation.Role.String()),
		Status:    domain.ToProjectInvitationStatus(invitation.Status.String()),
		ExpiresAt: invitation.ExpiresAt,
		CreatedAt: invitation.CreatedAt,
		UpdatedAt: invitation.UpdatedAt,
	}

	return result, nil
}

func (ProjectInvitationMapper) FromDomainToModel(invitation domain.ProjectInvitation) (*model.ProjectInvitations, error) {
	var result = &model.ProjectInvitations{
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
